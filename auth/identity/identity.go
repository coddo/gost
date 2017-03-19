package identity

import (
	"gost/auth/cookies"
	"gost/util/sliceutil"
)

// Identity represents the identity of the user the is in the current context
type Identity struct {
	Session      *cookies.Session
	User         *ApplicationUser
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

// HasAnyRole returns true if the current user has any of the specified roles
func (identity *Identity) HasAnyRole(roles []string) bool {
	if identity.User == nil {
		return false
	}

	return sliceutil.AreIntersected(identity.User.Roles, roles)
}

// New creates a new Identity based on a user session
func New(session *cookies.Session, user *ApplicationUser) *Identity {
	return &Identity{
		Session:      session,
		User:         user,
		isAuthorized: session != nil,
	}
}

// NewAnonymous creates a new anonymous Identity, with no user session defined
func NewAnonymous() *Identity {
	return &Identity{
		Session:      nil,
		User:         nil,
		isAuthorized: false,
	}
}
