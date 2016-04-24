package config

// Action struct represents an action that an endpoint has
type Action struct {
	Type           string `json:"type"`
	AllowAnonymous bool   `json:"allowAnonymous"`
	RequireAdmin   bool   `json:"requireAdmin"`
}

// Route entity
type Route struct {
	ID          string             `json:"id"`
	Endpoint    string             `json:"endpoint"`
	IsCacheable bool               `json:"isCacheable"`
	Actions     map[string]*Action `json:"actions"`
}

// Equal determines if the current Route is equal to another Route
func (route *Route) Equal(otherRoute *Route) bool {
	switch {
	case route.ID != otherRoute.ID:
		return false
	case route.Endpoint != otherRoute.Endpoint:
		return false
	case len(route.Actions) != len(otherRoute.Actions):
		return false
	default:
		for key, action := range route.Actions {
			if otherAction, found := otherRoute.Actions[key]; found {
				if !action.Equal(otherAction) {
					return false
				}
			} else {
				return false
			}
		}

		return true
	}
}

// Equal determines if the current Action is equal to another Action
func (action *Action) Equal(otherAction *Action) bool {
	switch {
	case action.Type != otherAction.Type:
		return false
	case action.AllowAnonymous != otherAction.AllowAnonymous:
		return false
	case action.RequireAdmin != otherAction.RequireAdmin:
		return false
	}

	return true
}
