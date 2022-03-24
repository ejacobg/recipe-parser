package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	findIngredientList(doc)
}

// See https://pkg.go.dev/golang.org/x/net/html#example-Parse
func findIngredientList(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "ul" {
		for _, a := range node.Attr {
			if a.Key == "class" {
				fmt.Println(a.Val)
				break
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		findIngredientList(c)
	}
}
