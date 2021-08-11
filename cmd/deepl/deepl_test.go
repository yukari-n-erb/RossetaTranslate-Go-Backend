package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeepLAPI(t *testing.T) {
	t.Run("returns DeepL API string When Get Method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/deepl", nil)
		response := httptest.NewRecorder()

		TranslateServer(response, request)

		got := response.Body.String()
		want := "DeepL"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("returns DeepL API string When Post Method with request body", func(t *testing.T) {
		// values := map[string]string{"name":"Hello"}
		values := Message{
			Text: "Do you want to play a game?",
		}
		fmt.Println(values)

		jsonString, _ := json.Marshal(values)
		fmt.Println(string(jsonString))

		request, _ := http.NewRequest(http.MethodPost, "/deepl", bytes.NewBuffer(jsonString))
		request.Header.Add("Content-Type", "application/json")
		response := httptest.NewRecorder()

		TranslateServer(response, request)

		got := response.Body.String()
		want := "\"ゲームをしたいですか？\""

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
