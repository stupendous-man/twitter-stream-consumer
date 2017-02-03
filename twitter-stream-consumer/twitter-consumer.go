package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"os"
)

func main() {

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	tweets, _, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{Count: 5})

	for _, tweet := range tweets {

		simpleTweet := SimpleTweet{}

		simpleTweet.Text = tweet.Text

		for _, url := range tweet.Entities.Urls {
			simpleTweet.DisplayUrl = url.DisplayURL
			simpleTweet.ExpandedUrl = url.ExpandedURL
			simpleTweet.Url = url.URL
		}

		fmt.Println(simpleTweet)
	}

	//TODO: Insert tweet Entity URLs in mongodb
}

type SimpleTweet struct {
	Text string
	DisplayUrl string
	ExpandedUrl string
	Url string
}
