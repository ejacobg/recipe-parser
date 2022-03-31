package parser

import "testing"

// assert function

// setup function
// read data from static file instead of making network calls
// grab recipe card and use that for all tests

func TestSimple(t *testing.T) {
	t.Run("slow-cooker-mashed-potatoes", func(t *testing.T) {

	})

	t.Run("olive-oil-mashed-potatoes", func(t *testing.T) {

	})

	t.Run("fluffy-garlic-herb-mashed-potatoes", func(t *testing.T) {

	})
}

// Ingredients may partially (need example) or fully be contained within a link
func TestIngredientContainsLink(t *testing.T) {
	// Fully contained
	t.Run("beef-cabbage-stir-fry", func(t *testing.T) {

	})

	// Fully contained
	t.Run("beef-taco-pasta", func(t *testing.T) {

	})
}

func TestMultipleIngredientLists(t *testing.T) {
	t.Run("beef-cabbage-stir-fry", func(t *testing.T) {

	})

	t.Run("chili-roasted-potatoes", func(t *testing.T) {
		
	})
}
