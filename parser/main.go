package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatalln("Error: too few arguments")
	}

	baseURL := "http://localhost:8080/"
	resp, err := http.Get(baseURL + os.Args[1] + ".html")
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
