package config

// Http methods
const (
	GetHTTPMethod    = "GET"
	PostHTTPMethod   = "POST"
	PutHTTPMethod    = "PUT"
	DeleteHTTPMethod = "DELETE"
)

// Route entity
type Route struct {
	ID       string            `json:"id"`
	Endpoint string            `json:"endpoint"`
	Actions  map[string]string `json:"actions"`
}

// Equal determines if the current Route is equal to another route
func (route *Route) Equal(otherRoute Route) bool {
	switch {
	case route.ID != otherRoute.ID:
		return false
	case route.Endpoint != otherRoute.Endpoint:
		return false
	case len(route.Actions) != len(otherRoute.Actions):
		return false
	default:
		for key, value := range route.Actions {
			if otherValue, found := otherRoute.Actions[key]; found {
				if value != otherValue {
					return false
				}
			} else {
				return false
			}
		}

		return true
	}
}
