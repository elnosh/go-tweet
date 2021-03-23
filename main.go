package main

import (
	"log"
	"net/http"

	"github.com/miguelhun/go-tweet/twitter"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body := twitter.ListenHashtag()
		log.Print(string(body))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
