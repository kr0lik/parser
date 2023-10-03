package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	*http.Client
}

func NewClient() *Client {
	return &Client{&http.Client{}}
}

func (c *Client) get(url string) (interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while create request: %v", err)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while do request: %v", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while read respose: %v", err)
	}

	var payload interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshal respose: %v", err)
	}

	return payload, nil
}
