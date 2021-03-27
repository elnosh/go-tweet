package twitter

import (
	"encoding/json"
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

type SearchTweetResponse struct {
	Data []SearchTweet `json:"data"`
	Meta Metadata
}

type SearchTweet struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	AuthorID  string `json:"author_id"`
	Text      string `json:"text"`
}

type Metadata struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
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

func (c *Client) ListenHashtag() (*SearchTweetResponse, error) {
	hashtag := config.GetHashtag()
	hashtagURL := fmt.Sprintf(c.BaseURL+"/tweets/search/recent?query=%s&tweet.fields=author_id,created_at", hashtag)

	req, err := http.NewRequest("GET", hashtagURL, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var response SearchTweetResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
