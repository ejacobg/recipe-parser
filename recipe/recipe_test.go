package recipe

import (
	"os"
	"testing"

	"github.com/ejacobg/recipe-parser/models"
	"github.com/ejacobg/recipe-parser/parser"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

// assert function
func assert(t testing.TB, got, want *models.Recipe) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

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
	}

	for _, test := range tests {
		rc := initTest(t, "./test-data/responses/"+test+".html")
		if rc == nil {
			t.Error("Error: couldn't find recipe card")
			return
		}
		t.Run(
			test, func(t *testing.T) {
				got, err := FromHTML(rc)
				if err != nil {
					t.Error("Error:", err)
					return
				}
				want, err := FromJSON("./test-data/solutions/" + test + ".json")
				if err != nil {
					t.Error("Error:", err)
					return
				}
				assert(t, got, want)
			},
		)
	}
}

// Ingredients may partially (need example) or fully be contained within a link
func TestIngredientContainsLink(t *testing.T) {
	tests := []string{
		"beef-cabbage-stir-fry", // fully contained
		"beef-taco-pasta",       // fully contained
	}

	for _, test := range tests {
		rc := initTest(t, "./test-data/responses/"+test+".html")
		if rc == nil {
			t.Error("Error: couldn't find recipe card")
			return
		}
		t.Run(
			test, func(t *testing.T) {
				got, err := FromHTML(rc)
				if err != nil {
					t.Error("Error:", err)
					return
				}
				want, err := FromJSON("./test-data/solutions/" + test + ".json")
				if err != nil {
					t.Error("Error:", err)
					return
				}
				assert(t, got, want)
			},
		)
	}
}

func TestMultipleIngredientLists(t *testing.T) {
	tests := []string{
		"beef-cabbage-stir-fry",
		"chili-roasted-potatoes",
		"fluffy-garlic-herb-mashed-potatoes",
	}

	for _, test := range tests {
		rc := initTest(t, "./test-data/responses/"+test+".html")
		if rc == nil {
			t.Error("Error: couldn't find recipe card")
			return
		}
		t.Run(
			test, func(t *testing.T) {
				got, err := FromHTML(rc)
				if err != nil {
					t.Error("Error:", err)
					return
				}
				want, err := FromJSON("./test-data/solutions/" + test + ".json")
				if err != nil {
					t.Error("Error:", err)
					return
				}
				assert(t, got, want)
			},
		)
	}
}
