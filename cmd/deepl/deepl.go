package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/DaikiYamakawa/deepl-go"
	"github.com/joho/godotenv"
)

type Message struct {
	Text string `json:"text"`
}

type TranslateText struct {
	Message string `json:"message"`
}

func TranslateServer(w http.ResponseWriter, r *http.Request) {
	translator := strings.TrimPrefix(r.URL.Path, "/")
	if translator == "deepl" {
		if r.Method == http.MethodGet {
			fmt.Fprint(w, "DeepL")
			return
		}
		if r.Method == http.MethodPost {
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			// Unmarshal
			var msg Message
			err = json.Unmarshal(b, &msg)
			fmt.Println(string(msg.Text))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			msgTranslated := DeeplApiTranslate(msg.Text)

			values := TranslateText{
				Message: msgTranslated,
			}

			output, err := json.Marshal(values)
			fmt.Println(string(output))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("content-type", "application/json")
			w.Write(output)
		}

	}

}

func DeeplApiTranslate(text string) string {
	err := godotenv.Load(fmt.Sprintf("../../%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		fmt.Println("Can't read .env file")
	}

	cli, err := deepl.New("https://api-free.deepl.com", nil)
	if err != nil {
		fmt.Printf("Failed to create client:\n   %+v\n", err)
	}
	translateResponse, err := cli.TranslateSentence(context.Background(), text, "EN", "JA")
	if err != nil {
		fmt.Printf("Failed to translate text:\n   %+v\n", err)
	} else {
		fmt.Printf("%+v\n", translateResponse.Translations[0].Text)
	}

	translatedText := translateResponse.Translations[0].Text

	return translatedText
}
