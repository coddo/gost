package cookies

import (
	"bytes"
	"errors"
	"gost/security"
	"gost/util"
	"io/ioutil"
	"os"

	"gopkg.in/mgo.v2/bson"
)

// FileCookieStore manages cookies using files on the local drive
type FileCookieStore struct {
	location string
}

// ReadCookie fetches a cookie from the cookie store
func (store *FileCookieStore) ReadCookie(key string) (*Session, error) {
	fileName := fileLocation(store.location, key)
	encryptedData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	jsonData, err := security.Decrypt(encryptedData)
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
	fileName := fileLocation(store.location, cookie.Token)
	jsonData, err := util.SerializeJSON(cookie)
	if err != nil {
		return err
	}

	encryptedData, err := security.Encrypt(jsonData)
	if err != nil {
		return err
	}

	err = addUserToken(store.location, cookie.UserID, cookie.Token)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, encryptedData, os.ModeDevice)
}

// DeleteCookie deletes a cookie from the cookie storage
func (store *FileCookieStore) DeleteCookie(cookie *Session) error {
	fileName := fileLocation(store.location, cookie.Token)

	err := removeUserToken(store.location, cookie.UserID, cookie.Token)
	if err != nil {
		return err
	}

	return os.Remove(fileName)
}

// GetAllUserCookies returns all the cookies that a certain user has
func (store *FileCookieStore) GetAllUserCookies(userID bson.ObjectId) ([]*Session, error) {
	userTokens, err := getUserTokens(store.location, userID)
	if err != nil {
		return nil, err
	}

	var cookies []*Session
	for _, token := range userTokens {
		session, err := store.ReadCookie(token)
		if err != nil {
			return nil, err
		}

		cookies = append(cookies, session)
	}

	return cookies, err
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
func fileLocation(storeLocation, key string) string {
	var buffer bytes.Buffer

	buffer.WriteString(storeLocation)
	buffer.WriteRune('/')
	buffer.WriteString(key)

	return buffer.String()
}

// NewFileCookieStore creates a new NewFileCookieStore pointer entity
func NewFileCookieStore(storeLocation string) *FileCookieStore {
	store := &FileCookieStore{location: storeLocation}

	return store
}

func addUserToken(storeLocation string, userID bson.ObjectId, token string) error {
	userTokens, err := getUserTokens(storeLocation, userID)
	if err != nil {
		return err
	}

	var isAlreadyAdded bool
	for _, userToken := range userTokens {
		if userToken == token {
			isAlreadyAdded = true
			break
		}
	}

	if !isAlreadyAdded {
		userTokens = append(userTokens, token)
	}

	return saveUserTokens(storeLocation, userID, userTokens)
}

func removeUserToken(storeLocation string, userID bson.ObjectId, token string) error {
	userTokens, err := getUserTokens(storeLocation, userID)
	if err != nil {
		return err
	}

	tokenIndex := -1
	for index, userToken := range userTokens {
		if userToken == token {
			tokenIndex = index
			break
		}
	}

	if tokenIndex == -1 {
		return errors.New("User token does not exist")
	}

	userTokens = append(userTokens[:tokenIndex], userTokens[tokenIndex+1:]...)

	return saveUserTokens(storeLocation, userID, userTokens)
}

func getUserTokens(storeLocation string, userID bson.ObjectId) ([]string, error) {
	userIndexFile := fileLocation(storeLocation, userID.Hex())

	fileContent, err := ioutil.ReadFile(userIndexFile)
	if err != nil {
		return nil, err
	}

	var userTokens []string
	err = util.DeserializeJSON(fileContent, userTokens)
	if err != nil {
		return nil, err
	}

	return userTokens, nil
}

func saveUserTokens(storeLocation string, userID bson.ObjectId, tokens []string) error {
	userIndexFile := fileLocation(storeLocation, userID.Hex())

	jsonData, err := util.SerializeJSON(tokens)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(userIndexFile, jsonData, os.ModeDevice)
	return err
}
