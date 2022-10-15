package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ejacobg/recipe-parser/models"
	"github.com/ejacobg/recipe-parser/recipe"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/html"
)

var coll *mongo.Collection

func Handler(w http.ResponseWriter, r *http.Request) {
	// Connect to MongoDB
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(
		context.TODO(), options.Client().ApplyURI(uri),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}()
	db := os.Getenv("DB_NAME")
	coll = client.Database(db).Collection("recipes")

	switch r.Method {
	case "GET":
		err = getRecipe(w, r)
	case "POST":
		err = postRecipe(w, r)
	case "PUT":
		err = putRecipe(w, r)
	case "DELETE":
		err = deleteRecipe(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

// Retrieves an entry from the database or parses it directly from the website.
// If the item isn't found in the database, then it will be parsed.
// To check if the item is in the database, I check the URL fields.
// I could also simply parse the site directly and determine which version to return,
// but I don't want to make more networks calls to the website than necessary.
func getRecipe(w http.ResponseWriter, r *http.Request) error {
	var (
		source string
		rcp    *models.Recipe
	)
	query := r.URL.Query()
	if names, ok := query["name"]; ok && len(names) > 0 {
		source = canonicalize(names[0])
	} else {
		http.Error(w, "Error: no name given", http.StatusNotFound)
		return nil
	}

	if _, src := query["src"]; src {
		rcp, err := recipeFromSource(source)
		if err != nil {
			return err
		}
		err = writeRecipe(w, rcp, http.StatusOK)
		if err != nil {
			return err
		}
		return nil
	}

	filter := bson.D{{"url", source}}
	err := coll.FindOne(context.TODO(), filter).Decode(rcp)
	if err == mongo.ErrNoDocuments {
		rcp, err = recipeFromSource(source)
		if err != nil {
			return err
		}
	}
	err = writeRecipe(w, rcp, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// Adds a new record to the database, failing if that item already exists.
// Use PUT to update saved recipes.
func postRecipe(w http.ResponseWriter, r *http.Request) error {
	var (
		source string
		rcp    *models.Recipe
	)
	query := r.URL.Query()
	if names, ok := query["name"]; ok && len(names) > 0 {
		source = canonicalize(names[0])
	} else {
		http.Error(w, "Error: no name given", http.StatusBadRequest)
		return nil
	}

	filter := bson.D{{"url", source}}
	result := coll.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		http.Error(w, "recipe already exists", http.StatusBadRequest)
		return nil
	}

	rcp, err := recipeFromSource(source)
	if err != nil {
		return err
	}

	// I can do a more thorough check now that the recipe has been parsed.
	filter = bson.D{{"id", rcp.ID}}
	opts := options.Replace().SetUpsert(true)
	_, err = coll.ReplaceOne(context.TODO(), filter, rcp, opts)
	if err != nil {
		return err
	}
	err = writeRecipe(w, rcp, http.StatusCreated)
	if err != nil {
		return err
	}
	return nil
}

// Updates an existing record in the database. WILL NOT create a new record, use POST.
func putRecipe(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		id  int
		rcp *models.Recipe
	)
	query := r.URL.Query()
	if ids, ok := query["id"]; ok && len(ids) > 0 {
		id, err = strconv.Atoi(ids[0])
		if err != nil {
			return err
		}
	} else {
		http.Error(w, "Error: no id given", http.StatusBadRequest)
		return nil
	}

	filter := bson.D{{"id", id}}
	err = coll.FindOne(context.TODO(), filter).Decode(rcp)
	if err != nil {
		http.Error(w, "id does not exist", http.StatusNotFound)
	}

	rcp, err = recipeFromSource(rcp.URL)
	if err != nil {
		return err
	}

	_, err = coll.UpdateOne(context.TODO(), filter, rcp)
	if err != nil {
		return err
	}
	err = writeRecipe(w, rcp, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

func deleteRecipe(w http.ResponseWriter, r *http.Request) (err error) {
	var id int
	query := r.URL.Query()
	if ids, ok := query["id"]; ok && len(ids) > 0 {
		id, err = strconv.Atoi(ids[0])
		if err != nil {
			return err
		}
	} else {
		http.Error(w, "Error: no id given", http.StatusBadRequest)
		return nil
	}

	filter := bson.D{{"id", id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// If the name is correct, then the canonicalized version should match the Recipe.URL field.
func canonicalize(name string) string {
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
func recipeFromSource(source string) (*models.Recipe, error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return recipe.FromHTML(doc)
}

func writeRecipe(w http.ResponseWriter, rcp *models.Recipe, status int) error {
	res, err := json.Marshal(*rcp)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return nil
}
