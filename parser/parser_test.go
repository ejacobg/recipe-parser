package parser

import (
	"os"
	"testing"

	"golang.org/x/net/html"
)

// assert function

// setup function
// read data from static file instead of making network calls
func parseFromFile(filename string) (*html.Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return html.Parse(file)
}

// grab recipe card and use that for all tests

func TestSimple(t *testing.T) {
	tests := []string{
		"slow-cooker-mashed-potatoes",
		"olive-oil-mashed-potatoes",
		"fluffy-garlic-herb-mashed-potatoes",
	}

	for _, test := range tests {

		t.Run(test, func(t *testing.T) {

		})
	}
}

// Ingredients may partially (need example) or fully be contained within a link
func TestIngredientContainsLink(t *testing.T) {
	tests := []string{
		"beef-cabbage-stir-fry", // fully contained
		"beef-taco-pasta",       // fully contained
	}

	for _, test := range tests {

		t.Run(test, func(t *testing.T) {

		})
	}
}

func TestMultipleIngredientLists(t *testing.T) {
	tests := []string{
		"beef-cabbage-stir-fry",
		"chili-roasted-potatoes",
	}

	for _, test := range tests {

		t.Run(test, func(t *testing.T) {

		})
	}
}
