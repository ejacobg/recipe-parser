package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ejacobg/recipe-parser/parser"
	"github.com/ejacobg/recipe-parser/recipe"

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

var readJSON = flag.Bool("json", false, "constructs a Recipe from a JSON file, then prints it")

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()

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

	rc := parser.FindRecipeCard(doc)
	if rc == nil {
		log.Fatalln("Error: couldn't find recipe card")
	}

	// ils := parser.FindIngredientLists(rc)
	// for _, il := range ils {
	// 	parser.PrintNode(il)
	// }

	r, err := recipe.FromHTML(rc)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	dbPath := "./database/"
	err = r.SaveAs(dbPath + os.Args[1])
	if err != nil {
		log.Fatalln("Error:", err)
	}

	fmt.Println("Recipe saved to ./database/" + args[0])
}
