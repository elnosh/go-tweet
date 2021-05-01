package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/miguelhun/go-tweet/config"
)

const (
	baseURL string = "https://api.twitter.com/2"
)

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

var (
	once          sync.Once
	TwitterClient *Client
	hashtagURL    = baseURL + "/tweets/search/recent?query=%s"
	streamURL     = baseURL + "/tweets/search/stream"
)

func newClient() *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func init() {
	if TwitterClient == nil {
		once.Do(
			func() {
				TwitterClient = newClient()
			})
	}
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	c.apiKey = config.GetTwitterKey()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type TweetResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type SearchTweetResponse struct {
	Data []TweetResponse `json:"data"`
	Meta Metadata        `json:"metadata"`
}

type Metadata struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

func (c *Client) ListenHashtag() (*SearchTweetResponse, error) {
	hashtag := config.GetHashtag()

	req, err := http.NewRequest("GET", fmt.Sprintf(hashtagURL, hashtag), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(req)
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

	var response SearchTweetResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type TweetStreamResponse struct {
	Data          TweetResponse `json:"data"`
	MatchingRules []RulesMatch  `json:"matching_rules"`
}

type RulesMatch struct {
	ID  int    `json:"id"`
	Tag string `json:"tag"`
}

func (c *Client) StreamTweets(input chan TweetStreamResponse) {
	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		log.Print(err)
	}

	res, err := c.sendRequest(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	jsonDecoder := json.NewDecoder(res.Body)

	for jsonDecoder.More() {

		var tweet TweetStreamResponse
		err := jsonDecoder.Decode(&tweet)
		if err != nil {
			log.Print(err)
		}

		input <- tweet
	}

}
