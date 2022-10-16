package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ejacobg/recipe-parser/api/utils"
	"github.com/ejacobg/recipe-parser/models"
)

// This route is functionally similar to /api/recipe, except it uses MongoDB's Data API.

var (
	baseURL = "https://data.mongodb-api.com/app/data-gicsu/endpoint/data/v1/action/"
	req     required
)

type dataAPIKey string

// Required by all request types.
type required struct {
	DataSource string `json:"dataSource"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

type id struct {
	ID string `json:"id"`
}

// Added an underscore because it conflicts with "net/url" from recipe.go.
type _url struct {
	URL string `json:"url"`
}

type _id struct {
	ID int `json:"_id"`
}

type findOne[T id | _url] struct {
	required
	Filter     *T   `json:"filter,omitempty"`
	Projection *_id `json:"projection,omitempty"`
}

func (*findOne[T]) action() string {
	return "findOne"
}

type replaceOne[T id | _url] struct {
	required
	Filter      T             `json:"filter"`
	Replacement models.Recipe `json:"replacement"`
	Upsert      bool          `json:"upsert,omitempty"`
}

func (*replaceOne[T]) action() string {
	return "replaceOne"
}

type deleteOne[T id | _url] struct {
	required
	Filter T `json:"filter"`
}

func (*deleteOne[T]) action() string {
	return "deleteOne"
}

type actioner interface {
	action() string
}

func sendRequest(ctx context.Context, a actioner) (*http.Response, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", baseURL+a.action(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Api-Key":      {ctx.Value(dataAPIKey("DATA_API_KEY")).(string)},
		"Accept":       {"application/json"},
		"Content-Type": {"application/json"},
	}
	return http.DefaultClient.Do(req)
}

type response struct {
	Document *models.Recipe `json:"document"`
}

func Data(w http.ResponseWriter, r *http.Request) {
	key := os.Getenv("DATA_API_KEY")
	ctx := context.WithValue(context.Background(), dataAPIKey("DATA_API_KEY"), key)
	r = r.Clone(ctx)

	db := os.Getenv("DB_NAME")
	req = required{"Cluster0", db, "recipes"}

	var err error
	switch r.Method {
	case "GET":
		err = getData(w, r)
	case "POST":
		err = postData(w, r)
	case "PUT":
		err = putData(w, r)
	case "DELETE":
		err = deleteData(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func getData(w http.ResponseWriter, r *http.Request) (err error) {
	var (
		source string
		rcp    = response{&models.Recipe{}}
	)
	query := r.URL.Query()
	if names, ok := query["name"]; ok && len(names) > 0 {
		source = utils.Canonicalize(names[0])
	} else {
		http.Error(w, "Error: no name given", http.StatusNotFound)
		return nil
	}

	if _, src := query["src"]; src {
		rcp.Document, err = utils.RecipeFromSource(source)
		if err != nil {
			return
		}
		return utils.WriteRecipe(w, rcp.Document, http.StatusOK)
	}

	res, err := sendRequest(r.Context(), &findOne[_url]{req, &_url{source}, &_id{0}})
	if err != nil {
		rcp.Document, err = utils.RecipeFromSource(source)
		if err != nil {
			return
		}
		return utils.WriteRecipe(w, rcp.Document, http.StatusOK)
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&rcp)
	if err != nil {
		return
	}
	if rcp.Document == nil {
		rcp.Document, err = utils.RecipeFromSource(source)
		if err != nil {
			return
		}
	}
	return utils.WriteRecipe(w, rcp.Document, http.StatusOK)
}

func postData(w http.ResponseWriter, r *http.Request) error {
	var (
		source string
		rcp    = response{&models.Recipe{}}
	)
	query := r.URL.Query()
	if names, ok := query["name"]; ok && len(names) > 0 {
		source = utils.Canonicalize(names[0])
	} else {
		http.Error(w, "Error: no name given", http.StatusNotFound)
		return nil
	}

	res, err := sendRequest(r.Context(), &findOne[_url]{req, &_url{source}, &_id{1}})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&rcp)
	if err != nil {
		return err
	}
	if rcp.Document != nil {
		http.Error(w, "recipe already exists", http.StatusBadRequest)
		return nil
	}

	rcp.Document, err = utils.RecipeFromSource(source)
	if err != nil {
		return err
	}

	res, err = sendRequest(
		r.Context(), &replaceOne[id]{
			req, id{rcp.Document.ID}, *rcp.Document,
			true,
		},
	)
	if err != nil {
		return err
	}
	return utils.WriteRecipe(w, rcp.Document, http.StatusOK)
}

func putData(w http.ResponseWriter, r *http.Request) error {
	var (
		ID  string
		rcp = response{&models.Recipe{}}
	)
	query := r.URL.Query()
	if ids, ok := query["id"]; ok && len(ids) > 0 {
		ID = ids[0]
	} else {
		http.Error(w, "Error: no id given", http.StatusBadRequest)
		return nil
	}

	res, err := sendRequest(r.Context(), &findOne[id]{req, &id{ID}, &_id{0}})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&rcp)
	if err != nil {
		return err
	}
	if rcp.Document == nil {
		http.Error(w, "id does not exist", http.StatusNotFound)
		return nil
	}

	rcp.Document, err = utils.RecipeFromSource(rcp.Document.URL)
	if err != nil {
		return err
	}

	res, err = sendRequest(
		r.Context(), &replaceOne[id]{
			req, id{rcp.Document.ID}, *rcp.Document,
			false,
		},
	)
	if err != nil {
		return err
	}
	return utils.WriteRecipe(w, rcp.Document, http.StatusOK)
}

func deleteData(w http.ResponseWriter, r *http.Request) error {
	var ID string
	query := r.URL.Query()
	if ids, ok := query["id"]; ok && len(ids) > 0 {
		ID = ids[0]
	} else {
		http.Error(w, "Error: no id given", http.StatusBadRequest)
		return nil
	}

	_, err := sendRequest(r.Context(), &deleteOne[id]{req, id{ID}})
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
