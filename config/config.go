package config

import "os"

func GetTwitterKey() string {
	return os.Getenv("TWITTER_KEY")
}

func GetHashtag() string {
	return os.Getenv("HASHTAG")
}
