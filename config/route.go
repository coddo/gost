package config

// Http methods
const (
    GETHttpMethod    = "GET"
    POSTHttpMethod   = "POST"
    PUTHttpMethod    = "PUT"
    DELETEHttpMethod = "DELETE"
)

// Route entity
type Route struct {
    Id       string            `json:"id"`
    Pattern  string            `json:"pattern"`
    Handlers map[string]string `json:"handlers"`
}

func (route *Route) Equal(otherRoute Route) bool {
    switch {
    case route.Id != otherRoute.Id:
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
