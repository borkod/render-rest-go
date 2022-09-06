package render

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	HttpClient   *http.Client
	ApiKey       string
	Host         string
	ServicesBase string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HttpClient:   http.DefaultClient,
		ApiKey:       apiKey,
		Host:         "https://api.render.com",
		ServicesBase: "/v1/services/",
	}
}
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusCreated {
		return body, err
	} else {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}
}
