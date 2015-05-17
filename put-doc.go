package gopherlastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PutDocumentRequest struct {
	Index string
	Type  string
	Id    string
	Doc   interface{}
}

type PutDocumentResponse struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	Id      string `json:"_id"`
	Version int    `json:"_version"`
	Created bool   `json:"created"`
}

func NewPutDocumentRequest(index string, docType string, id string, doc interface{}) *PutDocumentRequest {
	return &PutDocumentRequest{
		Index: index,
		Type:  docType,
		Id:    id,
		Doc:   doc,
	}
}

func (c *Client) PutDocument(putDocumentRequest *PutDocumentRequest) (*PutDocumentResponse, error) {

	req, err := c.buildPutDocumentHttpRequest(putDocumentRequest)
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

	var putResponse PutDocumentResponse
	err = json.Unmarshal(resBody, &putResponse)

	if err != nil {
		return nil, err
	}

	return &putResponse, nil
}

func (c *Client) buildPutDocumentHttpRequest(putDocumentRequest *PutDocumentRequest) (*http.Request, error) {
	// Since we support using URL as the ID, we need to use Opaque URL
	// so the http library doesn't un-encode the url-as-id;
	// therefore, we need to create our own request by hand
	body, err := json.Marshal(putDocumentRequest.Doc)
	if err != nil {
		return nil, err
	}

	bodyLength := int64(len(body))
	bodyReader := ioutil.NopCloser(bytes.NewReader(body))

	return &http.Request{
		Method: "PUT",
		Host:   c.Host, // takes precendence over URL.Host
		URL: &url.URL{
			Host:   c.Host, //ignored
			Scheme: "http",
			Opaque: c.buildDocIdPath(putDocumentRequest),
		},
		Body:          bodyReader,
		ContentLength: bodyLength,
	}, nil
}

func (c *Client) buildDocIdPath(putDocumentRequest *PutDocumentRequest) string {
	return fmt.Sprintf("/%s/%s/%s",
		putDocumentRequest.Index,
		putDocumentRequest.Type,
		putDocumentRequest.Id)
}
