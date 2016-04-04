package cookies

import (
	"gost/service"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DatabaseCookieStore manages cookies using a database
type DatabaseCookieStore struct {
	location string
}

// ReadCookie fetches a cookie from the cookie store
func (store *DatabaseCookieStore) ReadCookie(key string) (*Session, error) {
	session, collection := service.Connect(store.location)
	defer session.Close()

	cookie := Session{}
	err := collection.Find(bson.M{"token": key}).One(&cookie)

	return &cookie, err
}

// WriteCookie writes a cookie in the cookie store. If that cookie already exists,
// it is overwritten
func (store *DatabaseCookieStore) WriteCookie(cookie *Session) error {
	session, collection := service.Connect(store.location)
	defer session.Close()

	err := collection.UpdateId(cookie.ID, cookie)
	if err == mgo.ErrNotFound {
		err = collection.Insert(cookie)
	}

	return err
}

// DeleteCookie deletes a cookie from the cookie storage
func (store *DatabaseCookieStore) DeleteCookie(key string) error {
	session, collection := service.Connect(store.location)
	defer session.Close()

	err := collection.Remove(bson.M{"token": key})
	return err
}

// Init initializes the cookie store
func (store *DatabaseCookieStore) Init() {
	session, collection := service.Connect(store.location)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	if store.hasTokenIndex(collection) {
		return
	}

	index := mgo.Index{
		Key: []string{"$text:token"},
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// hasTokenIndex verifies if there is already an index created for the token field of a collection
func (store *DatabaseCookieStore) hasTokenIndex(collection *mgo.Collection) bool {
	indexes, err := collection.Indexes()
	if err != nil {
		panic(ErrInitializationFailed)
	}

	for _, index := range indexes {
		for _, key := range index.Key {
			if key == "token" {
				return true
			}
		}
	}

	return false
}

// NewDatabaseCookieStore creates a new DatabaseCookieStore pointer entity
func NewDatabaseCookieStore(storeLocation string) *DatabaseCookieStore {
	store := &DatabaseCookieStore{location: storeLocation}

	return store
}
