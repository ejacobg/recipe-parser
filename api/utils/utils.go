package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/ejacobg/recipe-parser/models"
	"github.com/ejacobg/recipe-parser/recipe"
	"golang.org/x/net/html"
)

// If the name is correct, then the canonicalized version should match the Recipe.URL field.
func Canonicalize(name string) string {
	// Convert from "%2D" to "-"
	esc, err := url.PathUnescape(name)
	if err != nil {
		// Ignore any failures
		esc = name
	}
	esc = strings.TrimSpace(esc)
	esc = strings.TrimSuffix(esc, "/")
	return "https://www.budgetbytes.com/" + strings.ToLower(esc) + "/"
}

// "source" should be a canonicalized name.
func RecipeFromSource(source string) (*models.Recipe, error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("non-OK HTTP status from budgetbytes.com: " + resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return recipe.FromHTML(doc)
}

func WriteRecipe(w http.ResponseWriter, rcp *models.Recipe, status int) error {
	res, err := json.Marshal(*rcp)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return nil
}
