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
	UserID     bson.ObjectId `bson:"userId,omitempty" json:"userId"`
	Token      string        `bson:"token,omitempty" json:"token"`
	ExpireTime time.Time     `bson:"expireTime,omitempty" json:"expireTime"`
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

// ResetExpireTime resets the expire time target of the session.
// This also triggers a Save() action, to update the cookie store
func (session *Session) ResetExpireTime() error {
	session.ExpireTime = util.NextDateFromNow(tokenExpireTime)

	return session.Save()
}

// NewSession generates a new Session pointer that contains the given userID and
// a unique token used as an identifier
func NewSession(userID bson.ObjectId) (*Session, error) {
	token, err := util.GenerateUUID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		UserID:     userID,
		Token:      token,
		ExpireTime: util.NextDateFromNow(tokenExpireTime),
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