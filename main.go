package main

import (
	"log"
	"net/http"

	"github.com/miguelhun/go-tweet/redis"
	"github.com/miguelhun/go-tweet/twitter"
)

func main() {
	http.HandleFunc("/search", getTweetSearch)
	http.HandleFunc("/stream", getTweetStream)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTweetSearch(w http.ResponseWriter, r *http.Request) {
	tweetResp, err := twitter.TwitterClient.ListenHashtag()
	if err != nil {
		log.Print(err)
	}
	for _, tweet := range tweetResp.Data {
		log.Printf("Tweet: %s\n", tweet.Text)
	}
}

func getTweetStream(w http.ResponseWriter, r *http.Request) {
	tweetsChan := make(chan twitter.TweetStreamResponse)

	go twitter.TwitterClient.StreamTweets(tweetsChan)

	go func() {
		for currentTweet := range tweetsChan {
			err := publishTweet(currentTweet.Data.Text)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func publishTweet(tweet string) error {
	err := redis.RedisClient.Publish("tweetChannel", tweet).Err()
	if err != nil {
		return err
	}

	return nil
}
