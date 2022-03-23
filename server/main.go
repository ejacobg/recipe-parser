package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// https://pkg.go.dev/net/http#example-FileServer
	fmt.Println("Serving on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./recipes"))))
}
