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

		fmt.Println(tweet.Text)

		for _, url := range tweet.Entities.Urls {
			fmt.Println("Display URL:", url.DisplayURL)
			fmt.Println("Expanded URL:", url.ExpandedURL)
			fmt.Println("URL:", url.URL)
		}

		fmt.Println()
	}

	//TODO: Insert tweet Entity URLs in mongodb
}
