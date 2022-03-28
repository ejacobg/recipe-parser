package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"parser"
	"recipe"

	"golang.org/x/net/html"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatalln("Error: too few arguments")
	}

	baseURL := "http://www.budgetbytes.com/"
	resp, err := http.Get(baseURL + os.Args[1])
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	rc := parser.FindRecipeCard(doc)
	if rc == nil {
		log.Fatalln("Error: couldn't find recipe card")
	}

	r, err := recipe.FromHTML(rc)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	json, err := r.ToJSON()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	fmt.Println(json)

}
