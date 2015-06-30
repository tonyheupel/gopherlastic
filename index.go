package gopherlastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type IndexAllocation struct {
	DisableAllocation string `json:"disable_allocation"` // Actually a string of bool
}

type IndexRouting struct {
	Allocation *IndexAllocation `json:"allocation"`
}

type Analyzer struct {
	Type     string `json:"type"`
	Language string `json:"language"`
}

type IndexAnalysis struct {
	Analyzer map[string]*Analyzer `json:"analyzer"`
}

type IndexVersion struct {
	Created string `json:"created"`
}

type IndexSettingsDetails struct {
	Routing          *IndexRouting  `json:"routing"`
	UUID             string         `json:"uuid"`
	NumberOfReplicas string         `json:"number_of_replicas"`
	Analysis         *IndexAnalysis `json:"analysis"`
	NumberOfShards   string         `json:"number_of_shards"`
	RefreshInterval  string         `json:"refresh_interval"`
	Version          *IndexVersion  `json:"version"`
}

type IndexSettings struct {
	Index *IndexSettingsDetails `json:"index"`
}

type Mapping struct {
	Properties *MappingProperties `json:"properties"`
}

type MappingProperties map[string]interface{}

type IndexDescription struct {
	Aliases  map[string]interface{} `json:"aliases"`
	Mappings map[string]*Mapping    `json:"mappings"`
	Settings *IndexSettings         `json:"settings"`
}
type GetIndexDescriptionResponse map[string]*IndexDescription

func (c *Client) GetIndexDescription(indexOrAliasName string) (*GetIndexDescriptionResponse, error) {

	req, err := c.buildGetIndexDescriptionHttpRequest(indexOrAliasName)
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

	var getResponse GetIndexDescriptionResponse
	err = json.Unmarshal(resBody, &getResponse)

	if err != nil {
		return nil, err
	}

	return &getResponse, nil
}

func (c *Client) buildGetIndexDescriptionHttpRequest(indexOrAliasName string) (*http.Request, error) {
	requestUrl, err := c.buildGetIndexDescriptionUrl(indexOrAliasName)

	if err != nil {
		return nil, err
	}

	return http.NewRequest("GET", requestUrl, nil)
}

func (c *Client) buildGetIndexDescriptionUrl(indexOrAliasName string) (string, error) {
	if indexOrAliasName == "" {
		return "", errors.New("Index or Alias name can not be blank")
	}

	return fmt.Sprintf("http://%s/%s", c.Host, indexOrAliasName), nil
}
