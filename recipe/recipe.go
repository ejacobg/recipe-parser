package recipe

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ejacobg/recipe-parser/models"
	"github.com/ejacobg/recipe-parser/parser"
	"golang.org/x/net/html/atom"

	"golang.org/x/net/html"
)

// FromHTML takes a document and attempts to build a recipe from it.
func FromHTML(doc *html.Node) (*models.Recipe, error) {
	recipeCard := parser.FindRecipeCard(doc)
	if recipeCard == nil {
		return nil, errors.New("couldn't find recipe card")
	}

	ingredientLists := parser.FindIngredientLists(recipeCard)
	if ingredientLists == nil {
		return nil, errors.New("couldn't find ingredients list(s)")
	}

	instructionsList := parser.FindInstructionsList(recipeCard)
	if instructionsList == nil {
		return nil, errors.New("instructions list does not exist")
	}

	return &models.Recipe{
		ID:           getID(recipeCard),
		Name:         getName(recipeCard),
		URL:          getURL(doc),
		Image:        getImage(recipeCard),
		Ingredients:  ingredientsFromLists(ingredientLists),
		Instructions: getInstructions(instructionsList),
	}, nil
}

// FromJSON reads from a JSON file and returns a recipe.
func FromJSON(path string) (*models.Recipe, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	recipe := &models.Recipe{}
	err = json.Unmarshal(file, &recipe)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

func getID(rc *html.Node) string {
	for _, a := range rc.Attr {
		if a.Key == "data-recipe-id" {
			return a.Val
		}
	}
	return ""
}

// Probably replace nil checks with panics because they shouldn't realistically happen
func getName(rc *html.Node) string {
	headerNode := parser.GetElementWithClass(rc, atom.H2, "wprm-recipe-name wprm-block-text-bold")
	if headerNode == nil {
		return "Error: headerNode not found"
	}
	textNode := headerNode.FirstChild
	if textNode == nil {
		return "Error: textNode not found"
	}
	return textNode.Data
}

func getURL(doc *html.Node) string {
	matcher := func(node *html.Node) bool {
		if node.Type == html.ElementNode && node.DataAtom == atom.Link {
			for _, a := range node.Attr {
				if a.Key == "rel" && a.Val == "canonical" {
					return true
				}
			}
		}
		return false
	}
	node := parser.FindNode(doc, matcher)
	if node == nil {
		return ""
	}
	for _, a := range node.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

func getImage(rc *html.Node) string {
	// The class list in the Elements tab has a different order than what is actually written in the raw HTML
	// Code from the HTTP response (line 999) looks like this: lazy lazy-hidden attachment-200x200 size-200x200
	// The rendered HTML uses this: lazy-hidden attachment-200x200 size-200x200
	imgNode := parser.GetElementWithClass(
		rc, atom.Img, "attachment-268x268 size-268x268 perfmatters-lazy",
	)
	if imgNode == nil {
		return "Error: imgNode not found"
	}
	for _, a := range imgNode.Attr {
		if a.Key == "data-pin-media" {
			return a.Val
		}
	}
	return "Error: could not find image link"
}

// Some recipes may have multiple ingredients lists
func ingredientsFromLists(lists []*html.Node) (ingredients []models.Ingredient) {
	for _, list := range lists {
		ingredients = append(ingredients, getIngredients(list)...)
	}
	return
}

// Assuming that the instructions list is parsed in order
func getInstructions(list *html.Node) []string {
	var instructions []string
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.DataAtom == atom.Li {
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

// Assumes ingredient list is passed
func getIngredients(list *html.Node) []models.Ingredient {
	classes := []string{
		"wprm-recipe-ingredient-amount",
		"wprm-recipe-ingredient-unit",
		"wprm-recipe-ingredient-name",
		"wprm-recipe-ingredient-notes wprm-recipe-ingredient-notes-normal",
	}
	var ingredients []models.Ingredient
	for li := list.FirstChild; li != nil; li = li.NextSibling {
		ingredient := models.Ingredient{}
		for index, class := range classes {
			spanNode := parser.GetElementWithClass(li, atom.Span, class)
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
