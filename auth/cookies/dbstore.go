package cookies

import (
	"gost/dal/service"

	"log"

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
func (store *DatabaseCookieStore) DeleteCookie(cookie *Session) error {
	session, collection := service.Connect(store.location)
	defer session.Close()

	err := collection.Remove(bson.M{"token": cookie.Token})
	return err
}

// GetAllUserCookies returns all the cookies that a certain user has
func (store *DatabaseCookieStore) GetAllUserCookies(userID bson.ObjectId) ([]*Session, error) {
	session, collection := service.Connect(store.location)
	defer session.Close()

	var userSessions []*Session
	err := collection.Find(bson.M{"userID": userID}).All(&userSessions)

	return userSessions, err
}

// DeleteAllUserCookies deletes all the cookies that a user has
func (store *DatabaseCookieStore) DeleteAllUserCookies(userID bson.ObjectId) error {
	session, collection := service.Connect(store.location)
	defer session.Close()

	_, err := collection.RemoveAll(bson.M{"userID": userID})

	return err
}

// ClearCookieStore clears the entire cookie store by deleting everything for all users
func (store *DatabaseCookieStore) ClearCookieStore() error {
	session, collection := service.Connect(store.location)
	defer session.Close()

	_, err := collection.RemoveAll(nil)

	return err
}

// Init initializes the cookie store and makes sure it is empty
func (store *DatabaseCookieStore) Init() {
	session, collection := service.Connect(store.location)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	index := mgo.Index{
		Key: []string{"$text:token"},
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		log.Fatalf("Database store indexes initialization error: %v\n", err)
	}

	err = store.ClearCookieStore()
	if err != nil {
		log.Fatalf("Database store initialization error for clearing previous cookies: %v\n", err)
	}
}

// NewDatabaseCookieStore creates a new DatabaseCookieStore pointer entity
func NewDatabaseCookieStore(storeLocation string) *DatabaseCookieStore {
	store := &DatabaseCookieStore{location: storeLocation}

	return store
}
