Scrapes recipes from budgetbyes.com and returns them as HTML or JSON. Includes other features for saving and querying recipes. Written in Go.

## Notes

Todo (unordered):

-   [ ] Make gallery of all the recipe pics
-   [ ] Gather all relevant info into data structure
-   [ ] Create API routes
-   [ ] Create webpage
-   [ ] Create DB of saved recipes
-   [ ] Create local/plaintext DB
-   [ ] Create queries/data visualizations (eg. show recipe vs ingredients)
    -   [ ] Add this to website

## Usage

When testing, remember to run the server (`cd server` then `go run .`).  
Only 1 recipe is active right now for the parser to use (`cd parser` then `go run . slow-cooker-mashed-potatoes`).

## Acknowledgements

https://github.com/poundifdef/plainoldrecipe
