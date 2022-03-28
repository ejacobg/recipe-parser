Scrapes recipes from budgetbyes.com and returns them as HTML or JSON. Includes other features for saving and querying recipes. Written in Go.

## Notes

Navigating to the "Print Recipe" link will bring you to a "minified" version of the recipe. This link contains the ID of the recipe, which might also be useful. The recipe ID is also found within the container div.

Todo (unordered):

-   [ ] Make gallery of all the recipe pics
-   [x] Gather all relevant info into data structure
-   [ ] Create API routes
-   [ ] Create webpage
-   [ ] Create DB of saved recipes
-   [x] Create local/plaintext DB
-   [ ] Create queries/data visualizations (eg. show recipe vs ingredients)
    -   [ ] Add this to website
-   [ ] Handle instructions that may have a nested list (find an example)
-   [x] Create and initialize recipe model
    -   [x] Test initialization of the recipe model
-   [ ] Create a "GetElement" or "GetElementByKeyValue" function

## Usage

When testing, remember to run the server (`cd server` then `go run .`).  
Only 1 recipe is active right now for the parser to use (`go run . slow-cooker-mashed-potatoes` in root).
Note: Right now the parser **has to be run in root** in order to access the `database` directory.

Some example recipes are included in the database for reference.

## Acknowledgements

https://github.com/poundifdef/plainoldrecipe
