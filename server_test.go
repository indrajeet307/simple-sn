package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexPage(t *testing.T) {
	t.Run("Test the index page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		IndexHandler(response, request)
		got := response.Body.String()
		want := HuddleGreeting

		if got != want {
			t.Errorf("got `%q`, expected `%s`", got, want)
		}
	})
}
