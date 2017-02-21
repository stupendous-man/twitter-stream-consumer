//TODO: Move to different package
package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"os/signal"
	"time"
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

		simpleTweet := SimpleTweet{}

		//Populate SimpleTweet struct
		simpleTweet.Text = tweet.Text

		for _, url := range tweet.Entities.Urls {
			simpleTweet.DisplayUrl = url.DisplayURL
			simpleTweet.ExpandedUrl = url.ExpandedURL
			simpleTweet.Url = url.URL
		}

		fmt.Println(simpleTweet)

		mongoInsert(simpleTweet)
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

//TODO: Move to separate go source file
func mongoInsert(simpleTweet SimpleTweet) {
	//Establish session with Mongo server
	session, err := mgo.Dial(os.Getenv("MONGO_PORT_27017_TCP_ADDR") + ":" + os.Getenv("MONGO_PORT_27017_TCP_PORT"))

	if err != nil {
		//Retry a few times before going into panic
		log.Print(err)
		for retries := 1; retries <= 5; retries++ {
			log.Printf("Retrying Mongo connection. Attempt %d...", retries)
			session, err = mgo.Dial(os.Getenv("MONGO_PORT_27017_TCP_ADDR") + ":" + os.Getenv("MONGO_PORT_27017_TCP_PORT"))

			if err != nil {
				time.Sleep(100 * time.Millisecond)
			} else {
				break
			}
		}
	}

	if err != nil {
		panic(err)
	}

	//Deferred function to close session
	defer session.Close()

	//Establish connection with database/collection
	c := session.DB("tweetsdb").C("tweets")

	//Insert document in database/collection
	err = c.Insert(&simpleTweet)

	if err != nil {
		log.Fatal(err)
	}
}

//SimpleTweet struct represents a processed tweet
type SimpleTweet struct {
	Text        string
	DisplayUrl  string
	ExpandedUrl string
	Url         string
}
