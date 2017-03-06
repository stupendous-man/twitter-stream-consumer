package mongo

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

func Insert(processedTweet ProcessedTweet) {
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
	err = c.Insert(&processedTweet)

	if err != nil {
		log.Fatal(err)
	}
}

//ProcessedTweet struct represents a processed tweet
type ProcessedTweet struct {
	Text        string
	DisplayUrl  string
	ExpandedUrl string
	Url         string
}
