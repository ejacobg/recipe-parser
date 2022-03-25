package main

import (
	"errors"
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

	list, err := findIngredientList(doc)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	printIngredientList(list)
}

// See https://pkg.go.dev/golang.org/x/net/html#example-Parse
func findIngredientList(node *html.Node) (*html.Node, error) {
	if node.Type == html.ElementNode && node.Data == "ul" {
		for _, a := range node.Attr {
			if a.Key == "class" && a.Val == "wprm-recipe-ingredients" {
				return node, nil
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		list, _ := findIngredientList(c)
		if list != nil {
			return list, nil
		}
	}
	return nil, errors.New("ingredient list does not exist")
}

func printIngredientList(list *html.Node) {
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		for child := li.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "span" {
				for _, a := range child.Attr {
					if a.Val == "wprm-recipe-ingredient-name" {
						// The first child should be a text node
						fmt.Println(child.FirstChild.Data)
					}
				}
			}
		}
	}
}
