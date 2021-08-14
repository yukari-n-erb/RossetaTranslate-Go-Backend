package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/translate"
	"github.com/DaikiYamakawa/deepl-go"
	"github.com/joho/godotenv"
	"golang.org/x/text/language"
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
			DeepLPost(w, r)
		}
	}
	if translator == "google" {
		if r.Method == http.MethodGet {
			fmt.Fprint(w, "Google")
			return
		}
		if r.Method == http.MethodPost {
			GooglePost(w, r)
		}

	}
}

func TextUnmarshal(w http.ResponseWriter, r *http.Request, msg *Message) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func TextMarshal(w http.ResponseWriter, r *http.Request, translatedText string) {
	values := TranslateText{
		Message: translatedText,
	}

	output, err := json.Marshal(values)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func DeepLPost(w http.ResponseWriter, r *http.Request) {
	var msg Message
	TextUnmarshal(w, r, &msg)

	fmt.Println("DeepL input:", msg.Text)
	msgTranslated := DeeplApiTranslate(msg.Text)
	fmt.Println("DeepL output:", msgTranslated)
	TextMarshal(w, r, msgTranslated)
}

func GooglePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Google")
	var msg Message
	TextUnmarshal(w, r, &msg)

	fmt.Println("Google input:", msg.Text)
	msgTranslated, _ := GoogleApiTranslate(msg.Text)
	fmt.Println("Google output:", msgTranslated)
	TextMarshal(w, r, msgTranslated)
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
	}

	translatedText := translateResponse.Translations[0].Text

	return translatedText
}

func GoogleApiTranslate(text string) (string, error) {
	err := godotenv.Load(fmt.Sprintf("../../%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		fmt.Println("Can't read .env file")
	}

	ctx := context.Background()

	lang, err := language.Parse("ja")
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}
