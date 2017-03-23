//TODO: Move to different package
package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/stupendous-man/twitter-stream-consumer/mongo-api"
	"log"
	"os"
	"os/signal"
)

func main() {

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	//Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()

	//TODO: Replace deprecated element
	//Set what kinds of stream input to process and configure output. In this case I'm only interested in processing tweets.
	demux.Tweet = func(tweet *twitter.Tweet) {
		ProcessTweet(tweet)
	}

	//demux.DM = func(dm *twitter.DirectMessage) {
	//	fmt.Println(dm.SenderID)
	//}
	//demux.Event = func(event *twitter.Event) {
	//	fmt.Printf("%#v\n", event)
	//}

	fmt.Println("Starting Stream...")

	//Set user params to establish stream
	userParams := &twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
		With:          "followings",
		Language:      []string{"en"},
	}

	stream, err := client.Streams.User(userParams)

	if err != nil {
		log.Fatal(err)
	}

	//Trap SIGINT to trigger shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	//TODO: Debug issue consuming stream within docker container
	//Receive twitter stream on a separate goroutine
	go func() {
		demux.HandleChan(stream.Messages)
	}()

	//Block main thread until SIGINT notification received
	for range signals {
		log.Println("Stopping Stream...")
		stream.Stop()
		break
	}

}

func ProcessTweet(tweet *twitter.Tweet) {

	processedTweet := mongo.ProcessedTweet{}

	//Populate ProcessedTweet struct
	processedTweet.Text = tweet.Text

	for _, url := range tweet.Entities.Urls {
		processedTweet.DisplayUrl = url.DisplayURL
		processedTweet.ExpandedUrl = url.ExpandedURL
		processedTweet.Url = url.URL
	}

	fmt.Println(processedTweet)

	mongo.Insert(processedTweet)

}