package twitter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/miguelhun/go-tweet/config"
)

const baseURL string = "https://api.twitter.com/2"

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) sendRequest(req *http.Request) ([]byte, error) {
	c.apiKey = config.GetTwitterKey()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%s", body))
	}
	return body, nil
}
