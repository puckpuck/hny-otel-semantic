package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HoneycombClient struct {
	baseUrl string
	client  *http.Client
}

type honeycombTransport struct {
	apiKey string
}

func (t *honeycombTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Honeycomb-Team", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}

func NewHoneycombClient(apiKey string) *HoneycombClient {
	return &HoneycombClient{
		baseUrl: "https://api.honeycomb.io",
		client:  &http.Client{Transport: &honeycombTransport{apiKey: apiKey}},
	}
}

func (c *HoneycombClient) ListAllDatasets() ([]HoneycombDataset, error) {

	url := c.baseUrl + "/1/datasets"
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var datasets []HoneycombDataset
	if err := json.Unmarshal(body, &datasets); err != nil {
		return nil, err
	}

	return datasets, nil
}

func (c *HoneycombClient) ListAllColumns(dataset HoneycombDataset) ([]HoneycombColumn, error) {

	url := c.baseUrl + "/1/columns/" + dataset.Name
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var columns []HoneycombColumn
	if err := json.Unmarshal(body, &columns); err != nil {
		return nil, err
	}

	return columns, nil
}

func (c *HoneycombClient) UpdateColumn(dataset HoneycombDataset, column HoneycombColumn) error {

	url := c.baseUrl + "/1/columns/" + dataset.Name + "/" + column.Id
	payload, err := json.Marshal(column)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error: status code: %s, while updating column: %s", resp.Status, column.KeyName)
	}

	return nil
}

type HoneycombDataset struct {
	ExpandJsonDepth     int    `json:"expand_json_depth"`
	RegularColumnsCount int    `json:"regular_columns_count"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Slug                string `json:"slug"`
	CreatedAt           string `json:"created_at"`
	LastWrittenAt       string `json:"last_written_at"`
}

type HoneycombColumn struct {
	Id          string `json:"id"`
	KeyName     string `json:"key_name"`
	Hidden      bool   `json:"hidden"`
	Description string `json:"description"`
	Type        string `json:"type"`
	LastWritten string `json:"last_written"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type HoneycombColumnUpdate struct {
	KeyName     string `json:"key_name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
}
