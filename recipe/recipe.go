package recipe

import (
	"errors"
	"models"
	"parser"

	"golang.org/x/net/html"
)

// Just return nil if it can't find the item?
// For best performance, let node point to the recipe card
func FromHTML(node *html.Node) (*models.Recipe, error) {
	ingredientLists := parser.FindIngredientLists(node)
	if ingredientLists == nil {
		return nil, errors.New("couldn't find ingredients list(s)")
	}

	instructionsList, err := parser.FindInstructionsList(node)
	if err != nil {
		return nil, err
	}

	return &models.Recipe{
		ID:           getID(node),
		Name:         getName(node),
		Image:        getImage(node),
		Ingredients:  IngredientsFromLists(ingredientLists),
		Instructions: getInstructions(instructionsList),
	}, nil
}

//
func getID(node *html.Node) string {
	rc := parser.FindRecipeCard(node)
	if rc == nil {
		// panic or return empty string
	}
	for _, a := range rc.Attr {
		if a.Key == "data-recipe-id" {
			return a.Val
		}
	}
	return ""
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
	// Code from the HTTP response (line 999) looks like this: lazy lazy-hidden attachment-200x200 size-200x200
	// The rendered HTML uses this: lazy-hidden attachment-200x200 size-200x200
	imgNode := parser.GetElementWithClass(node, "img", "lazy lazy-hidden attachment-200x200 size-200x200")
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
func getIngredients(list *html.Node) []models.Ingredient {
	classes := []string{
		"wprm-recipe-ingredient-amount",
		"wprm-recipe-ingredient-unit",
		"wprm-recipe-ingredient-name",
		"wprm-recipe-ingredient-notes wprm-recipe-ingredient-notes-normal",
	}
	ingredients := []models.Ingredient{}
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		ingredient := models.Ingredient{}
		for index, class := range classes {
			spanNode := parser.GetElementWithClass(li, "span", class)
			if spanNode == nil {
				// not all ingredients define all 4 classes
				continue
			}
			// Sometimes ingredients may be contained within links
			textNode := parser.GetTextNode(spanNode)
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

// Some recipes may have multiple ingredients lists
func IngredientsFromLists(lists []*html.Node) (ingredients []models.Ingredient) {
	for _, list := range lists {
		ingredients = append(ingredients, getIngredients(list)...)
	}
	return
}

// Assuming that the instructions list is parsed in order
func getInstructions(list *html.Node) []string {
	instructions := []string{}
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.Data == "li" {
			for _, a := range li.Attr {
				if a.Key == "class" && a.Val == "wprm-recipe-instruction" {
					textNode := parser.GetTextNode(li)
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
