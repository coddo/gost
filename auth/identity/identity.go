package identity

import "gost/auth/cookies"

// Identity represents the identity of the user the is in the current context
type Identity struct {
	Token       string
	Session     *cookies.Session
	isAnonymous bool
}

// IsAnonymous returns true if the current user is anonymous
func (identity *Identity) IsAnonymous() bool {
	return identity.isAnonymous
}

// IsAdmin returns true if the current authorized user is in the admin role
func (identity *Identity) IsAdmin() bool {
	return identity.Session.IsUserInRole(AccountTypeAdministrator)
}
