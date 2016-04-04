package identity

import "gost/auth/cookies"

// Identity represents the identity of the user the is in the current context
type Identity struct {
	Session      *cookies.Session
	isAuthorized bool
}

// IsAnonymous returns true if the current user is anonymous
func (identity *Identity) IsAnonymous() bool {
	return !identity.isAuthorized
}

// IsAuthorized returns true if the user is authorized
func (identity *Identity) IsAuthorized() bool {
	return identity.isAuthorized
}

// IsAdmin returns true if the current authorized user is in the admin role
func (identity *Identity) IsAdmin() bool {
	return identity.IsAuthorized() && identity.Session.IsUserInRole(AccountTypeAdministrator)
}

// New creates a new Identity based on a user session
func New(session *cookies.Session) *Identity {
	return &Identity{
		Session:      session,
		isAuthorized: session != nil,
	}
}

// NewAnonymous creates a new anonymous Identity, with no user session defined
func NewAnonymous() *Identity {
	return &Identity{
		Session:      nil,
		isAuthorized: false,
	}
}
