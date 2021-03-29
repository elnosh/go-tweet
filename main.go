package main

import (
	"log"
	"net/http"

	"github.com/miguelhun/go-tweet/twitter"
)

func main() {
	http.HandleFunc("/search", getTweetSearch)
	http.HandleFunc("/stream", getTweetStream)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTweetSearch(w http.ResponseWriter, r *http.Request) {
	client := twitter.NewClient()
	tweetResp, err := client.ListenHashtag()
	if err != nil {
		log.Print(err)
	}
	for _, tweet := range tweetResp.Data {
		log.Printf("Tweet: %s\n", tweet.Text)
	}
}

func getTweetStream(w http.ResponseWriter, r *http.Request) {
	client := twitter.NewClient()
	client.StreamTweets()
}
