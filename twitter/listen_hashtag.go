package twitter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/miguelhun/go-tweet/config"
)

var (
	url = "https://api.twitter.com/2/tweets/search/recent?query=%s&tweet.fields=author_id,created_at"
)

type TwitterResponse struct {
	Data []TwitterHashtagResponse `json:"data"`
	Meta Metadata
}
type TwitterHashtagResponse struct {
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	AuthorID  string `json:"author_id"`
	ID        string `json:"id"`
}

type Metadata struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

func ListenHashtag() []byte {
	client := http.Client{}

	hashtag := config.GetHashtag()
	urlHashtag := fmt.Sprintf(url, hashtag)

	req, err := http.NewRequest("GET", urlHashtag, nil)
	if err != nil {
		log.Fatal(err)
	}

	accessKey := config.GetTwitterKey()

	req.Header.Add("Authorization", accessKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}
