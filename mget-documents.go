package gopherlastic

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MGetDocumentsRequest struct {
	Documents []*GetDocumentRequest `json:"docs"`
}

type MGetDocumentsResponse struct {
	Documents []*GetDocumentResponse `json:"docs"`
}

func NewMGetDocumentsRequest(documentRequests []*GetDocumentRequest) *MGetDocumentsRequest {
	return &MGetDocumentsRequest{
		Documents: documentRequests,
	}
}

func (c *Client) MGetDocuments(mgetDocumentsRequest *MGetDocumentsRequest) (*MGetDocumentsResponse, error) {

	req, err := c.buildMGetDocumentHttpRequest(mgetDocumentsRequest)
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

	var mgetResponse MGetDocumentsResponse
	err = json.Unmarshal(resBody, &mgetResponse)

	if err != nil {
		return nil, err
	}

	return &mgetResponse, nil
}

func (c *Client) buildMGetDocumentHttpRequest(req *MGetDocumentsRequest) (*http.Request, error) {
	// TODO: Build post body
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	bodyLength := int64(len(body))
	bodyReader := ioutil.NopCloser(bytes.NewReader(body))

	return &http.Request{
		Method: "POST",
		Host:   c.Host, // takes precendence over URL.Host
		URL: &url.URL{
			Host:   c.Host, //ignored
			Scheme: "http",
			Path:   "_mget",
		},
		Body:          bodyReader,
		ContentLength: bodyLength,
	}, nil
}
