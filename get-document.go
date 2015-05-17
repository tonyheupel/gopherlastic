package gopherlastic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GetDocumentRequest struct {
	Index string
	Type  string
	Id    string
}

type GetDocumentResponse struct {
	Index   string           `json:"_index"`
	Type    string           `json:"_type"`
	Id      string           `json:"_id"`
	Version int              `json:"_version"`
	Source  *json.RawMessage `json:"_source"`
	Found   bool             `json:"found"`
}

func NewGetDocumentRequest(index string, docType string, id string) *GetDocumentRequest {
	return &GetDocumentRequest{
		Index: index,
		Type:  docType,
		Id:    id,
	}
}

func (c *Client) GetDocument(getDocumentRequest *GetDocumentRequest) (*GetDocumentResponse, error) {

	req, err := c.buildGetDocumentHttpRequest(getDocumentRequest)
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

	var getResponse GetDocumentResponse
	err = json.Unmarshal(resBody, &getResponse)

	if err != nil {
		return nil, err
	}

	return &getResponse, nil
}

func (c *Client) buildGetDocumentHttpRequest(req *GetDocumentRequest) (*http.Request, error) {
	// Since we support using URL as the ID, we need to use Opaque URL
	// so the http library doesn't un-encode the url-as-id;
	// therefore, we need to create our own request by hand
	return &http.Request{
		Method: "GET",
		Host:   c.Host, // takes precendence over URL.Host
		URL: &url.URL{
			Host:   c.Host, //ignored
			Scheme: "http",
			Opaque: buildDocIdPath(req.Index, req.Type, req.Id),
		},
	}, nil
}
