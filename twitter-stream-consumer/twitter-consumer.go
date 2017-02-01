package main

import (
	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
	"os"
	"fmt"
)

func main() {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	tweets, _, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{Count: 10})

	fmt.Printf("User's HOME TIMELINE:\n%+v\n", tweets)

	//TODO: Retrieve tweet Entities
	//TODO: Insert tweet Entities in mongodb
}