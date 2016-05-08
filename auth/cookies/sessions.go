package cookies

import (
	"errors"
	"gost/util"
	"gost/util/dateutil"
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
	ID          bson.ObjectId `bson:"_id" json:"-"`
	UserID      bson.ObjectId `bson:"userID,omitempty" json:"userID"`
	Token       string        `bson:"token,omitempty" json:"token"`
	AccountType int           `bson:"accountType,omitempty" json:"accountType"`
	ExpireTime  time.Time     `bson:"expireTime,omitempty" json:"-"`
	Client      *Client       `bson:"client,omitempty" json:"client"`
}

// Save saves the session in the cookie store
func (session *Session) Save() error {
	return cookieStore.WriteCookie(session)
}

// Delete deletes the session from the cookie store
func (session *Session) Delete() error {
	return cookieStore.DeleteCookie(session)
}

// IsExpired returns true if the session has expired
func (session *Session) IsExpired() bool {
	return dateutil.IsDateExpiredFromNow(session.ExpireTime)
}

// IsUserInRole verifies if the user with the current session has a specific role
func (session *Session) IsUserInRole(role int) bool {
	return session.AccountType == role
}

// ResetToken generates a new token and resets the expire time target of the session
// This also triggers a Save() action, to update the cookie store
func (session *Session) ResetToken() error {
	session.ExpireTime = dateutil.NextDateFromNow(tokenExpireTime)

	return session.Save()
}

// NewSession generates a new Session pointer that contains the given userID and
// a unique token used as an identifier
func NewSession(userID bson.ObjectId, accountType int, clientDetails *Client) (*Session, error) {
	token, err := util.GenerateUUID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:          bson.NewObjectId(),
		UserID:      userID,
		AccountType: accountType,
		Token:       token,
		ExpireTime:  dateutil.NextDateFromNow(tokenExpireTime),
		Client:      clientDetails,
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

		return nil, err
	}

	return session, nil
}

// GetUserSessions retrieves all the sessions that a user has
func GetUserSessions(userID bson.ObjectId) ([]*Session, error) {
	sessions, err := cookieStore.GetAllUserCookies(userID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(sessions); {
		if sessions[i].IsExpired() {
			err = sessions[i].Delete()
			if err == nil {
				err = ErrTokenExpired
			}

			sessions = append(sessions[:i], sessions[i+1:]...)
		} else {
			i++
		}
	}

	return sessions, nil
}

// ConfigureTokenExpireTime is used to change the default cofiguration of token expiration times
func ConfigureTokenExpireTime(expireTime time.Duration) {
	tokenExpireTime = expireTime
}
