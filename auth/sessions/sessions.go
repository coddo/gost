package auth

import (
	"bytes"
	"gost/util"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	expireTime          = 24 * 7 * time.Hour
	cookieStoreLocation = "cookies/"
)

// Session is a struct representing the session that a user has.
// Sessions are active since login until they expire or the user disconnects
type Session struct {
	UserID     bson.ObjectId
	Token      string
	ExpireTime time.Time
}

func GetSession(token string) *Session {
	return nil
}

func (session *Session) Save() error {
	return nil
}

func (session *Session) HasExpired() bool {
	return util.IsDateExpiredFromNow(session.ExpireTime)
}

func (session *Session) ResetExpireTime() {
	session.ExpireTime = util.NextDateFromNow(expireTime)
}

func getCoockieLocation(token string) string {
	var buffer bytes.Buffer

	buffer.WriteString(cookieStoreLocation)
	buffer.WriteString(token)

	return buffer.String()
}
