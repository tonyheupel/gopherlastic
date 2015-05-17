// The gopherlastic package provides convenient access to
// commonly used Elasticsearch features.
package gopherlastic

import (
	"fmt"
)

type Client struct {
	Host string
}

func NewClient(host string) *Client {
	return &Client{
		Host: host,
	}
}

func buildDocIdPath(index string, docType string, id string) string {
	return fmt.Sprintf("/%s/%s/%s", index, docType, id)
}
