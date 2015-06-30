// The golastic package provides convenient access to
// commonly used Elasticsearch features.
package gopherlastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// A Hit represents a single search result.
type Hit struct {
	Id     string                   `json:"_id"`
	Index  string                   `json:"_index"`
	Score  float64                  `json:"_score"`
	Source interface{}              `json:"_source"`
	Fields map[string][]interface{} `json:"fields"`
}

// Hits represents the summary of all search result hits.
type Hits struct {
	Total    int64   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

// Shards represents information about Shards that were queried
// as part of a search operation.
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// SearchResults represents the response to making a search
// against Elasticsearch.
type SearchResults struct {
	Took     int64  `json:"took"` // Time in ms
	TimedOut bool   `json:"timed_out"`
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
	Error    string `json:"error,omitempty"`
}

// SimpleSearchRequest represents a simplified representation of a keyword-based query
// structure that can be passed to Elasticsearch.
type SimpleSearchRequest struct {
	Skip     int64
	Count    int64
	Keywords string
	Index    string
	Type     string
}

// NewSimpleSearchRequest creates a new SimpleSearchRequest object based on the information passed in.
func NewSimpleSearchRequest(index string, docType string, keywords string, skip int64, count int64) *SimpleSearchRequest {
	return &SimpleSearchRequest{
		Index:    index,
		Type:     docType,
		Keywords: keywords,
		Skip:     skip,
		Count:    count,
	}
}

// SearchRequest is the type for generic JSON-based search requests to elasticsearch
type SearchRequest struct {
	Index string
	Type  string
	Body  string
}

// NewSearchRequest creates a new SearchRequest object based on the information passed in.
func NewSearchRequest(index string, docType string, body string) *SearchRequest {
	return &SearchRequest{
		Index: index,
		Type:  docType,
		Body:  body,
	}
}

// SimpleSearch performs a simple search against Elasticsearch.
// If keywords are provided, it matches all fields on the index using an "or" operator.
// If no keywords are provided, it returns all documents for the passed in index and type.
func (c *Client) SimpleSearch(req *SimpleSearchRequest) (*SearchResults, error) {
	// Search Using Raw json String if there are keywords passed in, otherwise do query.match_all = {}
	var countJson string
	if req.Count > 0 {
		countJson = fmt.Sprintf(", size: %d", req.Count)
	}

	var matchClause string
	if req.Keywords != "" {
		matchClause = fmt.Sprintf("\"query\": { multi_match: { _all: { query: \"%s\", operator: \"or\", type: \"cross_fields\", fields: [\"title\", \"description\", \"displayUrl\", \"metaKeywords\"] } } }", req.Keywords)
	} else {
		matchClause = "\"query\": { \"match_all\": {} }"
	}

	searchJson := fmt.Sprintf("{ from: %d%s, %s }", req.Skip, countJson, matchClause)

	url := fmt.Sprintf("http://%s/%s/%s/_search", c.Host, req.Index, req.Type)

	resp, err := http.Post(url, "application/json", strings.NewReader(searchJson))
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	var results SearchResults
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	return &results, err
}

func (c *Client) Search(searchRequest *SearchRequest) (*SearchResults, error) {
	req, err := c.buildSearchRequest(searchRequest)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil || (res.StatusCode >= 400 && res.StatusCode <= 599) {
		return nil, errors.New(string(resBody))
	}

	var searchResults SearchResults
	err = json.Unmarshal(resBody, &searchResults)

	if err != nil {
		return nil, err
	}

	if searchResults.Error != "" {
		return nil, errors.New(searchResults.Error)
	}

	return &searchResults, nil
}

func (c *Client) buildSearchRequest(req *SearchRequest) (*http.Request, error) {
	// Since we support using URL as the ID, we need to use Opaque URL
	// so the http library doesn't un-encode the url-as-id;
	// therefore, we need to create our own request by hand
	httpReq, err := http.NewRequest("POST",
		"http:"+buildSearchByIndexAndTypePath(c.Host, req.Index, req.Type),
		ioutil.NopCloser(strings.NewReader(req.Body)))
	if err != nil {
		return nil, err
	}
	return httpReq, nil
}

func buildSearchByIndexAndTypePath(host string, index string, docType string) string {
	var topicPart string

	if docType != "" {
		topicPart = docType + "/"
	}

	return fmt.Sprintf("//%s/%s/%s_search", host, index, topicPart)
}
