package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_appication_handler(t *testing.T){
	var registered = []struct {
		url  string
		expectedStatusCode int
	}{
		{"/", http.StatusOK},
		{"/something", http.StatusNotFound},
	}

	app := application{}
	routes := app.routes()

	test_server := httptest.NewServer(routes)

	defer test_server.Close()

	pathToTemplate = "./../../template/"

	for _, test := range registered {
		resp, err := test_server.Client().Get(test_server.URL+test.url)

		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != test.expectedStatusCode {
			t.Error("unexpeceted status code")
		}
	}
}