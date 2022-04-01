package recipe

import (
	"os"
	"parser"
	"testing"

	"golang.org/x/net/html"
)

// assert function

// setup function
func initTest(t testing.TB, name string) *html.Node {
	t.Helper()
	doc, err := parseFromFile(name)
	if err != nil {
		t.Error("Error parsing file:", err)
		return nil
	}
	rc := parser.FindRecipeCard(doc)
	return rc
}

// read data from static file instead of making network calls
func parseFromFile(name string) (*html.Node, error) {
	file, err := os.Open(name)
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
		rc := initTest(t, test)
		if rc == nil {
			t.Errorf("Error: couldn't find recipe card")
		}
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
		rc := initTest(t, test)
		if rc == nil {
			t.Errorf("Error: couldn't find recipe card")
		}

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
		rc := initTest(t, test)
		if rc == nil {
			t.Errorf("Error: couldn't find recipe card")
		}

		t.Run(test, func(t *testing.T) {

		})
	}
}
