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
	Pattern  string            `json:"pattern"`
	Handlers map[string]string `json:"handlers"`
}

// Equal determines if the current Route is equal to another route
func (route *Route) Equal(otherRoute Route) bool {
	switch {
	case route.ID != otherRoute.ID:
		return false
	case route.Pattern != otherRoute.Pattern:
		return false
	case len(route.Handlers) != len(otherRoute.Handlers):
		return false
	default:
		for key, value := range route.Handlers {
			if otherValue, found := otherRoute.Handlers[key]; found {
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
