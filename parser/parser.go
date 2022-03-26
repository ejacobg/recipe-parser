package parser

import (
	"errors"
	"fmt"

	"golang.org/x/net/html"
)

// See https://pkg.go.dev/golang.org/x/net/html#example-Parse
func FindIngredientList(node *html.Node) (*html.Node, error) {
	if node.Type == html.ElementNode && node.Data == "ul" {
		for _, a := range node.Attr {
			if a.Key == "class" && a.Val == "wprm-recipe-ingredients" {
				return node, nil
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		list, _ := FindIngredientList(c)
		if list != nil {
			return list, nil
		}
	}
	return nil, errors.New("ingredient list does not exist")
}

func PrintIngredientList(list *html.Node) {
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