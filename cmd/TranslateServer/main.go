package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(TranslateServer)
	log.Fatal(http.ListenAndServe(":3001", handler))
}
