package cookies

import (
	"bytes"
	"gost/util"
	"io/ioutil"
	"os"
)

// FileCookieStore manages cookies using files on the local drive
type FileCookieStore struct {
	location string
}

// ReadCookie fetches a cookie from the cookie store
func (store *FileCookieStore) ReadCookie(key string) (*Session, error) {
	fileName := store.fileLocation(key)
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

	return session, err
}

// WriteCookie writes a cookie in the cookie store. If that cookie already exists,
// it is overwritten
func (store *FileCookieStore) WriteCookie(cookie *Session) error {
	fileName := store.fileLocation(cookie.Token)
	jsonData, err := util.SerializeJSON(cookie)
	if err != nil {
		return err
	}

	encodedData := util.Encode(jsonData)

	return ioutil.WriteFile(fileName, encodedData, os.ModeDevice)
}

// DeleteCookie deletes a cookie from the cookie storage
func (store *FileCookieStore) DeleteCookie(key string) error {
	fileName := store.fileLocation(key)

	return os.Remove(fileName)
}

// Init initializes the cookie store
func (store *FileCookieStore) Init() {
	if _, err := os.Stat(store.location); os.IsNotExist(err) {
		err = os.Mkdir(store.location, os.ModeDevice)

		if err != nil {
			panic(ErrInitializationFailed)
		}
	}
}

// cookieLocation computes and returns the location of a certain cookie
func (store *FileCookieStore) fileLocation(key string) string {
	var buffer bytes.Buffer

	buffer.WriteString(store.location)
	buffer.WriteRune('/')
	buffer.WriteString(key)

	return buffer.String()
}

// NewFileCookieStore creates a new NewFileCookieStore pointer entity
func NewFileCookieStore(storeLocation string) *FileCookieStore {
	store := &FileCookieStore{location: storeLocation}

	return store
}