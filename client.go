// The gopherlastic package provides convenient access to
// commonly used Elasticsearch features.
package gopherlastic

type Client struct {
	Host string
}

func NewClient(host string) *Client {
	return &Client{
		Host: host,
	}
}
