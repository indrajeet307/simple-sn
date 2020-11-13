package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexPage(t *testing.T) {
	t.Run("Test the index page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		IndexHandler(response, request)
		rj := Response{}

		json.Unmarshal(response.Body.Bytes(), &rj)

		want := HuddleGreeting

		if rj.Message != want {
			t.Errorf("got `%q`, expected `%s`", rj.Message, want)
		}
	})
}
