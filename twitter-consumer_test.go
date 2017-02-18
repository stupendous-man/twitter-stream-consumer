//TODO: Move to different package
package main

import (
	"encoding/json"
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

	//TODO:Try unmarshalling into Tweet struct from go-twitter/twitter
	mockTwitterResponse := TestTweet{}
	json.Unmarshal([]byte(mockServerResponseBody), &mockTwitterResponse)

	//TODO: Move SimpleTweet instantiation logic to func in twitter-consumer.go
	simpleTweet := TestSimpleTweet{}

	//Populate SimpleTweet struct
	simpleTweet.Text = mockTwitterResponse.Text

	for _, url := range mockTwitterResponse.Entities.Urls {
		simpleTweet.DisplayUrl = url.DisplayURL
		simpleTweet.ExpandedUrl = url.ExpandedURL
		simpleTweet.Url = url.URL
	}

	if simpleTweet.Text != "Test tweet text..." {
		t.Errorf("Expected \"Test tweet text...\", but instead got %s", simpleTweet.Text)
	}

	if simpleTweet.DisplayUrl != "xyz.co/123abC" {
		t.Errorf("Expected \"xyz.co/123abC\" but instead got \"%s\"", simpleTweet.DisplayUrl)
	}

	if simpleTweet.ExpandedUrl != "http://xyz.co/123abC" {
		t.Errorf("Expected \"http://xyz.co/123abC\" but instead got \"%s\"", simpleTweet.ExpandedUrl)
	}

	if simpleTweet.Url != "https://t.co/123aBc" {
		t.Errorf("Expected \"https://t.co/123aBc\" but instead got \"%s\"", simpleTweet.Url)
	}

}

//Simplified version of Tweet struct defined in go-twitter/twitter for test purposes
type TestTweet struct {
	Text     string        `json:"text"`
	Entities *TestEntities `json:"entities"`
}

//Simplified version of Entities struct defined in go-twitter/twitter for test purposes
type TestEntities struct {
	Urls []TestURLEntity `json:"urls"`
}

//Simplified version of URLEntity struct defined in go-twitter/twitter for test purposes
type TestURLEntity struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	URL         string `json:"url"`
}

//TestSimpleTweet struct represents a processed tweet
type TestSimpleTweet struct {
	Text        string
	DisplayUrl  string
	ExpandedUrl string
	Url         string
}
