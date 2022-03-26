package recipe

import (
	"errors"

	"golang.org/x/net/html"
)

type Ingredient struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
	Name   string `json:"name"`
	Notes  string `json:"notes"`
}

type Recipe struct {
	Name         string       `json:"name"`
	Image        string       `json:"image"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
}

// Just return nil if it can't find the item?
func FromHTML(doc *html.Node) (*Recipe, error) {
	ingredientList, err := findIngredientList(doc)
	if err != nil {
		return nil, err
	}

	instructionsList, err := findInstructionsList(doc)
	if err != nil {
		return nil, err
	}

	return &Recipe{
		getName(doc),
		getImage(doc),
		getIngredients(ingredientList),
		getInstructions(instructionsList),
	}, nil
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

func findInstructionsList(node *html.Node) (*html.Node, error) {

}

func getName(node *html.Node) string {

}

func getImage(node *html.Node) string {

}

func getIngredients(node *html.Node) []Ingredient {

}

func getInstructions(node *html.Node) []string {

}
