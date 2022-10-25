Scrapes recipes from budgetbytes.com and returns them as JSON. Includes other features for 
saving and querying recipes. Written in Go.

## Usage

The API endpoint for all requests can be found here:

```
https://recipe-parser-ejacobg.vercel.app/api/{data|recipe}
```

Both routes are functionally equivalent, however it is recommended to use the `/data` route for 
better performance. The API supports the following requests:

| Action | Parameters      | Description                                                                                                                                                 |
|--------|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GET    | name (`string`) | Required. Queries database for information on the given recipe. If not found, will parse from budgetbytes.com.                                              |
|        | src (`boolean`) | Optional. Skips the database query and parses directly from budgetbytes.com.                                                                                |
| POST   | name (`string`) | Required. Adds and returns the named recipe to the database. Fails if the recipe already exists.                                                            |
| PUT    | id (`string`)   | Required. Uses the "id" field of the recipe. Parses budgetbytes.com and updates the database entry. Returns the updated recipe. Fails if the id is unknown. |
| DELETE | id (`string`)   | Required. Deletes the recipe from the database.                                                                                                             |

The value for the `name` parameter is taken from the recipe's URL. For example, if we wanted to 
query for this recipe: https://www.budgetbytes.com/slow-cooker-mashed-potatoes/, then the `name` 
parameter would be `slow-cooker-mashed-potatoes`. For convenience, this bookmarklet will allow 
you to extract the recipe name from the URL:

```
javascript:(()=>{navigator.clipboard.writeText(document.location.pathname.split('/', 2)[1])})()
```

With the name, you can make a GET or POST request to the service:

```javascript
fetch(`https://recipe-parser-ejacobg.vercel.app/api/{data|recipe}`, { method: "GET|POST" })
    .then((response) => response.json())
    .then((data) => {/* ... */})
```

You can update the code as you see fit.

### Example

All successful responses (except for DELETE) look similar to this:

```json
{
  "id": "30990",
  "name": "Slow Cooker Mashed Potatoes",
  "url": "https://www.budgetbytes.com/slow-cooker-mashed-potatoes/",
  "image": "https://www.budgetbytes.com/wp-content/uploads/2015/12/Slow-Cooker-Mashed-Potatoes-scoop.jpg",
  "ingredients": [
    {
      "amount": "3",
      "unit": "lbs.",
      "name": "russet potatoes",
      "notes": "($1.80)"
    },
    {
      "amount": "1.5",
      "unit": "cups",
      "name": "chicken broth",
      "notes": "($0.20)"
    },
    {
      "amount": "2",
      "unit": "cloves",
      "name": "garlic, minced",
      "notes": "($0.16)"
    },
    {
      "amount": "1/4",
      "unit": "tsp",
      "name": "Freshly cracked black pepper",
      "notes": "($0.05)"
    },
    {
      "amount": "4",
      "unit": "oz.",
      "name": "cream cheese",
      "notes": "($0.40)"
    },
    {
      "amount": "1/2",
      "unit": "cup",
      "name": "milk",
      "notes": "($0.25)"
    },
    {
      "amount": "1",
      "unit": "Tbsp",
      "name": "butter",
      "notes": "($0.13)"
    }
  ],
  "instructions": [
    "Wash and peel the potatoes, then dice them into one-inch cubes. Rinse the diced potatoes with cool water in a colander to remove the excess starch.",
    "Add the cubed potatoes, minced garlic, chicken broth, and some freshly cracked pepper to the slow cooker. Stir briefly to distribute the garlic and pepper.",
    "Place a lid on the slow cooker and cook on high for three hours, or until the potatoes are fork tender. You can test the tenderness by lifting the lid just long enough to pierce the potatoes with a fork.",
    "Take the lid off the slow cooker and add the cream cheese, milk, and butter. Stir to combine the ingredients and mash the potatoes. For an extra smooth mashed potato, use a hand mixer to briefly whip the potatoes until smooth.",
    "Taste the potatoes and add salt or pepper if needed. Serve immediately, or switch the slow cooker to the \"warm\" setting until ready to serve."
  ]
}
```

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
-   [ ] Create webpage
-   [ ] Update main function to launch a server
-   [ ] Dockerize app

Nice-to-have:

-   [ ] Make gallery of all the recipe pics
-   [ ] Create queries/data visualizations (eg. show recipe vs ingredients)
    -   [ ] Add this to website
-   [ ] Handle instructions that may have a nested list (find an example)
-   [ ] Create "GetElement**s**" functions that return a collection rather than just the first one
-   [ ] Include the recipe link in the recipe model for convenience
-   [ ] Use `html.Render` instead of your `PrintNode` function
-   [ ] Allow batch processing from file

## Acknowledgements

https://github.com/poundifdef/plainoldrecipe
