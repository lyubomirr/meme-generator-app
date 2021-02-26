package main

import (
	"github.com/lyubomirr/meme-generator/presentation"
	"log"
	"net/http"
)

func main() {
	router := presentation.GetRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
