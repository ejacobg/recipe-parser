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


func getName(doc *html.Node) string {

}

func getImage(doc *html.Node) string {

}

func getIngredients(list *html.Node) []Ingredient {

}

func getInstructions(list *html.Node) []string {

}
