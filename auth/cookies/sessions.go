package cookies

import (
	"errors"
	"gost/util"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	defaultTokenExpireTime = 24 * 7 * time.Hour
)

var (
	tokenExpireTime = defaultTokenExpireTime
)

// Errors generated during session handling
var (
	ErrTokenExpired = errors.New("The session token has expired")
)

// Session is a struct representing the session that a user has.
// Sessions are active since login until they expire or the user disconnects
type Session struct {
	UserID      bson.ObjectId `bson:"userId,omitempty" json:"userId"`
	Token       string        `bson:"token,omitempty" json:"token"`
	AccountType int           `bson:"accountType,omitempty" json:"accountType"`
	ExpireTime  time.Time     `bson:"expireTime,omitempty" json:"expireTime"`
	Client      *Client       `bson:"client,omitempty" json:"client"`
}

// Client struct contains information regarding the client that has made the http request
type Client struct {
	IPAddress string  `bson:"ipAddress,omitempty" json:"ipAddress"`
	Browser   string  `bson:"browser,omitempty" json:"browser"`
	OS        string  `bson:"os,omitempty" json:"os"`
	Country   string  `bson:"country,omitempty" json:"country"`
	State     string  `bson:"state,omitempty" json:"state"`
	City      string  `bson:"city,omitempty" json:"city"`
	Latitude  float64 `bson:"latitude,omitempty" json:"latitude"`
	Longitude float64 `bson:"longitude,omitempty" json:"longitude"`
}

// Save saves the session in the cookie store
func (session *Session) Save() error {
	return cookieStore.WriteCookie(session)
}

// Delete deletes the session from the cookie store
func (session *Session) Delete() error {
	return cookieStore.DeleteCookie(session.Token)
}

// IsExpired returns true if the session has expired
func (session *Session) IsExpired() bool {
	return util.IsDateExpiredFromNow(session.ExpireTime)
}

// IsUserInRole verifies if the user with the current session has a specific role
func (session *Session) IsUserInRole(role int) bool {
	return session.AccountType == role
}

// ResetToken generates a new token and resets the expire time target of the session
// This also triggers a Save() action, to update the cookie store
func (session *Session) ResetToken() error {
	token, err := util.GenerateUUID()
	if err != nil {
		return err
	}

	session.Token = token
	session.ExpireTime = util.NextDateFromNow(tokenExpireTime)

	return session.Save()
}

// NewSession generates a new Session pointer that contains the given userID and
// a unique token used as an identifier
func NewSession(userID bson.ObjectId, client *Client) (*Session, error) {
	token, err := util.GenerateUUID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		UserID:     userID,
		Token:      token,
		ExpireTime: util.NextDateFromNow(tokenExpireTime),
		Client:     client,
	}

	return session, nil
}

// GetSession retrieves a session from the cookie store
func GetSession(token string) (*Session, error) {
	session, err := cookieStore.ReadCookie(token)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		err = session.Delete()
		if err == nil {
			err = ErrTokenExpired
		}

		return nil, ErrTokenExpired
	}

	return session, nil
}

// ConfigureTokenExpireTime is used to change the default cofiguration of token expiration times
func ConfigureTokenExpireTime(expireTime time.Duration) {
	tokenExpireTime = expireTime
}
