package recipe

import (
	"parser"

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
// For best performance, let node point to the recipe card
func FromHTML(node *html.Node) (*Recipe, error) {
	ingredientList, err := parser.FindIngredientList(node)
	if err != nil {
		return nil, err
	}

	instructionsList, err := parser.FindInstructionsList(node)
	if err != nil {
		return nil, err
	}

	return &Recipe{
		getName(node),
		getImage(node),
		getIngredients(ingredientList),
		getInstructions(instructionsList),
	}, nil
}

// Probably replace nil checks with panics because they shouldn't realistically happen
func getName(node *html.Node) string {
	headerNode := parser.GetElementWithClass(node, "h2", "wprm-recipe-name wprm-block-text-bold")
	if headerNode == nil {
		return "Error: headerNode not found"
	}
	textNode := headerNode.FirstChild
	if textNode == nil {
		return "Error: textNode not found"
	}
	return textNode.Data
}

func getImage(node *html.Node) string {
	// The class list in the Elements tab has a different order than what is actually written in the raw HTML
	imgNode := parser.GetElementWithClass(node, "img", "lazy-hidden attachment-200x200 size-200x200")
	if imgNode == nil {
		return "Error: imgNode not found"
	}
	for _, a := range imgNode.Attr {
		if a.Key == "data-pin-media" {
			return a.Val
		}
	}
	// Maybe just return empty string for this
	return "Error: could not find image link"
}

// Assumes ingredient list is passed
func getIngredients(list *html.Node) []Ingredient {
	classes := []string{
		"wprm-recipe-ingredient-amount",
		"wprm-recipe-ingredient-unit",
		"wprm-recipe-ingredient-name",
		"wprm-recipe-ingredient-notes wprm-recipe-ingredient-notes-normal",
	}
	ingredients := []Ingredient{}
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		ingredient := Ingredient{}
		for index, class := range classes {
			spanNode := parser.GetElementWithClass(li, "span", class)
			if spanNode == nil {
				// panic
			}
			textNode := spanNode.FirstChild
			switch index {
			case 0:
				ingredient.Amount = textNode.Data
			case 1:
				ingredient.Unit = textNode.Data
			case 2:
				ingredient.Name = textNode.Data
			case 3:
				ingredient.Notes = textNode.Data
			}
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients
}

// Assuming that the instructions list is parsed in order
func getInstructions(list *html.Node) []string {
	instructions := []string{}
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.Data == "li" {
			for _, a := range li.Attr {
				if a.Key == "class" && a.Val == "wprm-recipe-instruction" {
					// Cheating here
					// Create a "GetElement" or "GetElementByKeyValue" function
					textNode := li.FirstChild.FirstChild.FirstChild
					if textNode == nil {
						// panic
					}
					instructions = append(instructions, textNode.Data)
				}
			}
		}
	}
	return instructions
}
