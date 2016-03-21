package auth

import (
	"bytes"
	"errors"
	"gost/util"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	defaultTokenExpireTime = 24 * 7 * time.Hour
	cookieStoreLocation    = "cookies/"
)

var (
	tokenExpireTime = defaultTokenExpireTime
)

// Errors generated by session handling
var (
	ErrTokenExpired = errors.New("The session token has expired")
)

// Session is a struct representing the session that a user has.
// Sessions are active since login until they expire or the user disconnects
type Session struct {
	UserID     bson.ObjectId
	Token      string
	ExpireTime time.Time
}

// Save saves the session in the cookie store
func (session *Session) Save() error {
	jsonData, err := util.SerializeJSON(session)
	if err != nil {
		return err
	}

	encodedData := util.Encode(jsonData)
	fileName := getCoockieLocation(session.Token)

	return ioutil.WriteFile(fileName, encodedData, os.ModeDevice)
}

// Delete deletes the session from the cookie store
func (session *Session) Delete() error {
	fileName := getCoockieLocation(session.Token)

	return os.Remove(fileName)
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
	fileName := getCoockieLocation(token)
	encodedData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	jsonData, err := util.Decode(encodedData)
	if err != nil {
		return nil, err
	}

	var session *Session
	err = util.DeserializeJSON(jsonData, session)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		session.Delete()
		return nil, ErrTokenExpired
	}

	return session, session.ResetExpireTime()
}

// SetTokenExpireTime is used to change the default cofiguration of token expiration times
func SetTokenExpireTime(expireTime time.Duration) {
	tokenExpireTime = expireTime
}

func getCoockieLocation(token string) string {
	var buffer bytes.Buffer

	buffer.WriteString(cookieStoreLocation)
	buffer.WriteString(token)

	return buffer.String()
}