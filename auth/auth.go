package auth

import (
	"errors"
	"gost/auth/cookies"
	"gost/auth/identity"
	"gost/security"
	"gost/util"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// The keys that are used in the request header to authorize the user
const (
	AuthorizationHeader = "Authorization"
	AuthorizationScheme = "GHOST-TOKEN"
)

// Errors generated by the auth package
var (
	ErrInvalidScheme           = errors.New("The used authorization scheme is invalid or not supported")
	ErrInvalidToken            = errors.New("The given token is expired or invalid")
	ErrInvalidUser             = errors.New("There is no application user with the given ID")
	ErrDeactivatedUser         = errors.New("The current user account is deactivated")
	ErrInexistentClientDetails = errors.New("Missing client details. Cannot create authorization for anonymous client")

	errAnonymousUser = errors.New("The user has no identity")
)

// GenerateUserAuth generates a new gost-token, saves it in the database and returns it to the client
func GenerateUserAuth(userID bson.ObjectId, client *cookies.Client) (string, error) {
	if client == nil {
		return ErrInexistentClientDetails.Error(), ErrInexistentClientDetails
	}

	if !identity.IsUserExistent(userID) {
		return ErrInvalidUser.Error(), ErrInvalidUser
	}

	session, err := cookies.NewSession(userID, client)
	if err != nil {
		return err.Error(), err
	}

	err = session.Save()
	if err != nil {
		return err.Error(), err
	}

	ghostToken, err := generateGhostToken(session)

	return ghostToken, err
}

// Authorize tries to authorize an existing gostToken
func Authorize(httpHeader http.Header) (*identity.Identity, error) {
	ghostToken, err := extractGostToken(httpHeader)
	if err != nil {
		if err == errAnonymousUser {
			return identity.NewAnonymous(), nil
		}

		return nil, err
	}

	encryptedToken, err := util.Decode([]byte(ghostToken))
	if err != nil {
		return nil, err
	}

	jsonToken, err := security.Decrypt(encryptedToken)
	if err != nil {
		return nil, err
	}

	cookie := new(cookies.Session)
	err = util.DeserializeJSON(jsonToken, cookie)
	if err != nil {
		return nil, err
	}

	dbCookie, err := cookies.GetSession(cookie.Token)
	if err != nil || dbCookie == nil {
		return nil, err
	}

	if !identity.IsUserActivated(dbCookie.UserID) {
		return nil, ErrDeactivatedUser
	}

	go dbCookie.ResetToken()

	return identity.New(dbCookie), nil
}

func generateGhostToken(session *cookies.Session) (string, error) {
	jsonToken, err := util.SerializeJSON(session)
	if err != nil {
		return err.Error(), err
	}

	encryptedToken, err := security.Encrypt(jsonToken)
	if err != nil {
		return err.Error(), err
	}

	ghostToken := util.Encode(encryptedToken)

	return string(ghostToken), nil
}

func extractGostToken(httpHeader http.Header) (string, error) {
	var gostToken string

	if gostToken = httpHeader.Get(AuthorizationHeader); len(gostToken) == 0 {
		return errAnonymousUser.Error(), errAnonymousUser
	}

	if !strings.Contains(gostToken, AuthorizationScheme) {
		return ErrInvalidScheme.Error(), ErrInvalidScheme
	}

	gostTokenValue := strings.TrimPrefix(gostToken, AuthorizationScheme)
	gostTokenValue = strings.TrimSpace(gostTokenValue)

	if len(gostTokenValue) == 0 {
		return ErrInvalidToken.Error(), ErrInvalidToken
	}

	return gostTokenValue, nil
}
