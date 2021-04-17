package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis"
	"github.com/miguelhun/go-tweet/twitter"
)

var (
	once     sync.Once
	client   *twitter.Client
	rdClient *redis.Client
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

func rdsClient() *redis.Client {
	rdClient = redis.NewClient(&redis.Options{})
	// ask about this return and in getClient
	return rdClient
}

func main() {
	getClient()
	rdsClient()

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

	go func() {
		for currentTweet := range tweetsChan {
			err := publishTweet(currentTweet.Data.Text, rdClient)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
}

func publishTweet(tweet string, redisClient *redis.Client) error {
	err := redisClient.Publish("tweetChannel", tweet).Err()
	if err != nil {
		return err
	}

	return nil
}
