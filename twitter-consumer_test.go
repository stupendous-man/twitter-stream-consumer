//TODO: Move to different package
package main

import (
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/stupendous-man/twitter-stream-consumer/mongo-api"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTweetConsumption(t *testing.T) {
	handler := http.NotFound

	//Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		handler(rw, r)
	}))

	defer ts.Close()

	handler = func(rw http.ResponseWriter, r *http.Request) {
		//Create mock response
		resp := []byte(`{"text" : "Test tweet text...", "entities" : {"urls" : [ {"display_url" : "xyz.co/123abC", "expanded_url" : "http://xyz.co/123abC", "url" : "https://t.co/123aBc"} ]}}`)

		//Reply with mock response
		rw.Write([]byte(resp))
	}

	mockServerResponse, err := http.Get(ts.URL)

	if err != nil {
		log.Fatal(err)
	}

	mockServerResponseBody, err := ioutil.ReadAll(mockServerResponse.Body)
	mockServerResponse.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	//Unmarshal into Tweet struct from go-twitter/twitter
	mockTwitterResponse := twitter.Tweet{}
	json.Unmarshal([]byte(mockServerResponseBody), &mockTwitterResponse)

	processedTweet := mongo.ProcessedTweet{}

	//Populate SimpleTweet struct
	processedTweet.Text = mockTwitterResponse.Text

	for _, url := range mockTwitterResponse.Entities.Urls {
		processedTweet.DisplayUrl = url.DisplayURL
		processedTweet.ExpandedUrl = url.ExpandedURL
		processedTweet.Url = url.URL
	}

	if processedTweet.Text != "Test tweet text..." {
		t.Errorf("Expected \"Test tweet text...\", but instead got %s", processedTweet.Text)
	}

	if processedTweet.DisplayUrl != "xyz.co/123abC" {
		t.Errorf("Expected \"xyz.co/123abC\" but instead got \"%s\"", processedTweet.DisplayUrl)
	}

	if processedTweet.ExpandedUrl != "http://xyz.co/123abC" {
		t.Errorf("Expected \"http://xyz.co/123abC\" but instead got \"%s\"", processedTweet.ExpandedUrl)
	}

	if processedTweet.Url != "https://t.co/123aBc" {
		t.Errorf("Expected \"https://t.co/123aBc\" but instead got \"%s\"", processedTweet.Url)
	}

}
