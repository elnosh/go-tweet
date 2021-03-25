package main

import (
	"log"
	"net/http"

	"github.com/miguelhun/go-tweet/twitter"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client := twitter.NewClient()
		tweetResp, err := client.ListenHashtag()
		if err != nil {
			log.Print(err)
		}
		for _, tweet := range tweetResp.Data {
			log.Printf("Tweet: %s\n", tweet.Text)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
