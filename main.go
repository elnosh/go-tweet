package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/miguelhun/go-tweet/twitter"
)

var (
	once   sync.Once
	client *twitter.Client
)

func getClient() *twitter.Client {
	if client == nil {
		once.Do(
			func() {
				client = twitter.NewClient()
			})
	}
	return client
}

func main() {
	getClient()
	http.HandleFunc("/search", getTweetSearch)
	http.HandleFunc("/stream", getTweetStream)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTweetSearch(w http.ResponseWriter, r *http.Request) {
	tweetResp, err := client.ListenHashtag()
	if err != nil {
		log.Print(err)
	}
	for _, tweet := range tweetResp.Data {
		log.Printf("Tweet: %s\n", tweet.Text)
	}
}

func getTweetStream(w http.ResponseWriter, r *http.Request) {
	tweetsChan := make(chan twitter.TweetStreamResponse)

	go client.StreamTweets(tweetsChan)

	for currentTweet := range tweetsChan {
		log.Printf("tweet: %s", currentTweet.Data.Text)
	}

}
