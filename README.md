# RossetaTranslate-Go-Backend
This is a backend server for translating text by DeepL API or Google Translate API.

## Setup
1. clone this repogitory
2. make .env file
3. write DEEPL_API_KEY=<your DeepL API key> to .env file
4. write GOOGLE_APPLICATION_CREDENTIALS=<your Google Translate API service account json filepath> to .env file
3. cd cmd/TranslateServer
4. go build
5. go run ./main.go ./TranslateServer.go

## How can use
1. Create a request data with JSON data
{
    "text" : "<original text>"
}

2. post JSON data to following address
- use DeepL API
localhost:3000/deepl

- use Google API
localhost:3000/google
