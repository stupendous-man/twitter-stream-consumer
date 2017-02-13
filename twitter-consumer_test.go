package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
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
		resp := `{
			"Tweet" : {
				"Entities" : {
					"Urls" : [ {
						"DisplayURL" : "xyz.co/123abC",
						"ExpandedURL" : "http://xyz.co/123abC",
						"URL" : "https://t.co/123aBc",
					} ]
				}
			}
		}`

		//Reply with mock response
		rw.Write([]byte(resp))
	}

//	ms, err := MockService(ts.URL)
}

//func MockService(url string) (Service, error) {
//	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
//	if err != nil {
//		return nil, err
//	}
//	return &service{elasticClient: client}, nil
//}