package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/miguelhun/go-tweet/config"
)

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
