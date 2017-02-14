package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"
	"log"
	"io/ioutil"
)

func TestTwitterStreamConsumption(t *testing.T) {
	handler := http.NotFound

	//Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		handler(rw, r)
	}))

	defer ts.Close()

	handler = func(rw http.ResponseWriter, r *http.Request) {
		//Create mock response
		resp := []byte(`{"entities" : {"urls" : [ {"display_url" : "xyz.co/123abC", "expanded_url" : "http://xyz.co/123abC", "url" : "https://t.co/123aBc"} ]}}`)

		//Reply with mock response
		rw.Write([]byte(resp))
	}

	mockServerResponse, err := http.Get(ts.URL)

	if err != nil {
		log.Fatal(err)
	}

	MockTwitterResponse, err := ioutil.ReadAll(mockServerResponse.Body)
	mockServerResponse.Body.Close()

	if err != nil {

		log.Fatal(err)
	}

	fmt.Printf("%s", MockTwitterResponse)

	//TODO: Continue unit test...
}

//Simplified version of Tweet struct defined in go-twitter/twitter for test purposes
type Tweet struct {
	Entities *Entities `json:"entities"`
}

//Simplified version of Entities struct defined in go-twitter/twitter for test purposes
type Entities struct {
	Urls []URLEntity `json:"urls"`
}

//Simplified version of URLEntity struct defined in go-twitter/twitter for test purposes
type URLEntity struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	URL         string `json:"url"`
}
