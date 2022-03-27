package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"parser"

	"golang.org/x/net/html"
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

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	list, err := parser.FindIngredientList(doc)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	parser.PrintIngredientList(list)

	rc := parser.FindRecipeCard(doc)
	if rc == nil {
		log.Fatalln("Error: couldn't find recipe card")
	}
	for _, a := range rc.Attr {
		fmt.Println(a.Key, a.Val)
	}
}
