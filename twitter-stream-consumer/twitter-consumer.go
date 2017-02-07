package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"os"
	"os/signal"
	"gopkg.in/mgo.v2"
	"log"
	"time"
	"fmt"
	"syscall"
)

func main() {

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	//TODO: Retrieve twitter stream in separate go routine
	//tweets, _, _ := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{Count: 5})
	//
	//for _, tweet := range tweets {
	//
	//	simpleTweet := SimpleTweet{}
	//
	//	//Populate SimpleTweet struct
	//	simpleTweet.Text = tweet.Text
	//
	//	for _, url := range tweet.Entities.Urls {
	//		simpleTweet.DisplayUrl = url.DisplayURL
	//		simpleTweet.ExpandedUrl = url.ExpandedURL
	//		simpleTweet.Url = url.URL
	//	}
	//
	//	//Insert SimpleTweet in MongoDB
	//	mongoInsert(simpleTweet)
	//}

	//TODO: Clean this up a bit
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")

	// USER (quick test: auth'd user likes a tweet -> event)
	 userParams := &twitter.StreamUserParams{
	 	StallWarnings: twitter.Bool(true),
	 	With:          "followings",
	 	Language:      []string{"en"},
	 }
	 stream, err := client.Streams.User(userParams)
	 if err != nil {
	 	log.Fatal(err)
	 }

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	//TODO: Insert tweet events in Mongo
	//TODO: Handle SIGINT

	fmt.Println("Stopping Stream...")
	stream.Stop()
}

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

type SimpleTweet struct {
	Text string
	DisplayUrl string
	ExpandedUrl string
	Url string
}
