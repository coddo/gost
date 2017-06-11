package cookies

var (
	unknownClientDetails = &Client{
		Address: "unknown",
		Type:    "unknown",
		Name:    "unknown",
		Version: "unknown",
		OS:      "unknown",
	}
)

// Client struct contains information regarding the client that has made the http request
type Client struct {
	Address        string `bson:"address,omitempty" json:"address"`
	Type           string `bson:"type,omitempty" json:"type"`
	Name           string `bson:"name,omitempty" json:"name"`
	Version        string `bson:"version,omitempty" json:"version"`
	OS             string `bson:"os,omitempty" json:"os"`
	IsMobileDevice bool   `bson:"isMobileDevice,omitempty" json:"isMobileDevice"`
}

// UnknownClientDetails gets the default client details, which have 'unknown' values
func UnknownClientDetails() *Client {
	return unknownClientDetails
}

// Equal check if two client details are the same
func (client *Client) Equal(other *Client) bool {
	switch {
	case client.Address != other.Address:
		return false
	case client.Type != other.Type:
		return false
	case client.Name != other.Name:
		return false
	case client.Version != other.Version:
		return false
	case client.OS != other.OS:
		return false
	case client.IsMobileDevice != other.IsMobileDevice:
		return false
	default:
		return true
	}
}
