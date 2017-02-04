package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"os"
	"gopkg.in/mgo.v2"
	"log"
	"time"
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

func mongoInsert(simpleTweet SimpleTweet) {
	//Establish session with Mongo server
	session, err := mgo.Dial(os.Getenv("MONGO_PORT_27017_TCP_ADDR") + ":" + os.Getenv("MONGO_PORT_27017_TCP_PORT"))

	if err != nil {
		//Retry a few times before going into panic
		log.Print(err)
		for retries := 1; retries <= 5; retries++ {
			log.Printf("Retrying Mongo connection. Attempt %d...", retries)
			session, err = mgo.Dial("localhost")

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
	c := session.DB("goctupusdb").C("events")

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
