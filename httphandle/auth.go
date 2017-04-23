package httphandle

import (
	"gost/auth"
	"gost/auth/identity"
	"net/http"
)

func authorize(req *http.Request, allowAnonymous bool, roles []string) (*identity.Identity, bool, int) {
	identity, err := auth.Authorize(req.Header)
	if err != nil {
		return nil, false, http.StatusUnauthorized
	}

	if allowAnonymous {
		return identity, true, http.StatusOK
	}

	if identity.IsAnonymous() {
		return nil, false, http.StatusUnauthorized
	}

	if !identity.HasAnyRole(roles...) {
		return nil, false, http.StatusForbidden
	}
}
