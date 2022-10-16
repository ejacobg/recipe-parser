Scrapes recipes from budgetbyes.com and returns them as HTML or JSON. Includes other features for saving and querying recipes. Written in Go.

## Notes

Navigating to the "Print Recipe" link will bring you to a "minified" version of the recipe. This link contains the ID of the recipe, which might also be useful. The recipe ID is also found within the container div.

The "Print Recipe" page can be seen by going to `budgetbytes.com/wprm_print/<recipe_id>`

Good-to-have:

-   [x] Gather all relevant info into data structure
-   [x] Create API routes
-   [x] Create DB of saved recipes
-   [x] Create local/plaintext DB
-   [x] Create and initialize recipe model
    -   [x] Test initialization of the recipe model
-   [x] Create a "GetElement" or "GetElementByKeyValue" function
-   [x] Fix multiple ingredients issue (chili-roasted-potatoes is broken, but not others)
-   [x] Use `atom.Atom` to check for tag names instead of strings
-   [x] Refactor parser to use collector/matcher pattern (https://gist.github.
    com/Xeoncross/8bbb84bc4bf540bd907f79ee17c4e1fc)
-   [ ] Write test suite for parser

Nice-to-have:

-   [ ] Create webpage
-   [ ] Make gallery of all the recipe pics
-   [ ] Create queries/data visualizations (eg. show recipe vs ingredients)
    -   [ ] Add this to website
-   [ ] Handle instructions that may have a nested list (find an example)
-   [ ] Create "GetElement**s**" functions that return a collection rather than just the first one
-   [ ] Include the recipe link in the recipe model for convenience
-   [ ] Use `html.Render` instead of your `PrintNode` function
-   [ ] Allow batch processing from file

## Usage

When testing, remember to run the server (`cd server` then `go run .`).  
Only 1 recipe is active right now for the parser to use (`go run . slow-cooker-mashed-potatoes` in root).

Some example recipes are included in the database for reference.

## Acknowledgements

https://github.com/poundifdef/plainoldrecipe
