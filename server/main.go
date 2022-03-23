package main

import (
	"log"
	"net/http"
)

func main() {
	// https://pkg.go.dev/net/http#example-FileServer
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./recipes"))))
}
