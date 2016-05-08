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
