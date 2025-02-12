package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ejacobg/recipe-parser/recipe"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/net/html"
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage of recipe-parser:")
		fmt.Fprintln(w, "recipe-parser <recipe-name>")
		fmt.Fprintln(w, "recipe-parser -json <recipe-name>.json")
		fmt.Fprintln(w, "Obtain the recipe name from the budgetbytes.com URL.")
		flag.PrintDefaults()
	}
}

var (
	readJSON = flag.Bool("json", false, "constructs a Recipe from a JSON file, then prints it")
	dbPath   = flag.String("dbpath", "./database/", "path to save JSON output")
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()

	mongodb()

	if len(args) < 1 {
		log.Fatalln("Error: too few arguments")
	}

	if *readJSON {
		r, err := recipe.FromJSON(args[0])
		if err != nil {
			log.Fatalln("Error:", err)
		}
		fmt.Println(r)
		os.Exit(0)
	}

	baseURL := "http://www.budgetbytes.com/"
	resp, err := http.Get(baseURL + args[0])
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	r, err := recipe.FromHTML(doc)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	err = r.SaveAs(*dbPath + os.Args[1])
	if err != nil {
		log.Fatalln("Error:", err)
	}

	fmt.Println("Recipe saved to " + *dbPath + args[0])
}

func mongodb() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Println("MONGODB_URI variable must be set")
		return
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("mongo_dev").Collection("recipes")
	r, err := recipe.FromJSON("./database/slow-cooker-mashed-potatoes.json")
	if err != nil {
		log.Println("Could not read recipe from file")
	}

	filter := bson.D{{"id", r.ID}}
	opts := options.Replace().SetUpsert(true)
	result, err := coll.ReplaceOne(context.TODO(), filter, r, opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(
		"Matched:", result.MatchedCount, "Modified:", result.MatchedCount, "Upserted:",
		result.UpsertedCount,
	)
	return
}
