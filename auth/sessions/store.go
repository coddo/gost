// Package sessions uses the DatabaseCookieStore type as its default cookie store
// for managing user sessions
package sessions

import (
	"bytes"
	"gost/util"
	"io/ioutil"
	"os"
)

const (
	defaultCookieStoreLocation = "cookies"
)

// CookieStore is an interface used for describing entities that can manage
// the location of user sessions (cookies) and performs operations such as
// reading or writing cookie from/to that location
type CookieStore interface {
	Location() string
	ReadCookie(key string) ([]byte, error)
	WriteCookie(key string, data []byte) error
	DeleteCookie(key string) error
	cookieLocation(key string) string
}

var defaultCookieStore = &DatabaseCookieStore{location: defaultCookieStoreLocation}
var cookieStore CookieStore = defaultCookieStore

// DatabaseCookieStore manages cookies using a database
type DatabaseCookieStore struct {
	location string
}

// Location gets the full path/connection which contains the location of the cookie store data
func (store *DatabaseCookieStore) Location() string {
	return store.location
}

// ReadCookie fetches a cookie from the cookie store
func (store *DatabaseCookieStore) ReadCookie(key string) ([]byte, error) {
	fileName := cookieStore.cookieLocation(key)
	encodedData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return util.Decode(encodedData)
}

// WriteCookie writes a cookie in the cookie store. If that cookie already exists,
// it is overwritten
func (store *DatabaseCookieStore) WriteCookie(key string, data []byte) error {
	fileName := cookieStore.cookieLocation(key)
	encodedData := util.Encode(data)

	return ioutil.WriteFile(fileName, encodedData, os.ModeDevice)
}

// DeleteCookie deletes a cookie from the cookie storage
func (store *DatabaseCookieStore) DeleteCookie(key string) error {
	fileName := cookieStore.cookieLocation(key)

	return os.Remove(fileName)
}

// cookieLocation computes and returns the location of a certain cookie
func (store *DatabaseCookieStore) cookieLocation(key string) string {
	var buffer bytes.Buffer

	buffer.WriteString(cookieStore.Location())
	buffer.WriteRune('/')
	buffer.WriteString(key)

	return buffer.String()
}

// NewDatabaseCookieStore creates a new DatabaseCookieStore pointer entity
func NewDatabaseCookieStore(storeLocation string) *DatabaseCookieStore {
	store := &DatabaseCookieStore{location: storeLocation}

	return store
}

// SetCookieStore sets the cookie store that will be used by the system.
// If this method is not called, the default cookie store will be used
func SetCookieStore(store CookieStore) {
	cookieStore = store
}
