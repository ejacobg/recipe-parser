// Package parser contains all DOM traversal functions used by other packages.
package parser

import (
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// See https://pkg.go.dev/golang.org/x/net/html#example-Parse

// FindRecipeCard returns the node representing the enclosing <div> of the recipe.
// All the relevant recipe information is rooted at this node.
func FindRecipeCard(node *html.Node) *html.Node {
	return GetElementWithClass(node, atom.Div, "wprm-recipe-container")
}

// FindIngredientLists returns all nodes representing a list of ingredients.
// Some recipes group their ingredients under different headers (eg. "Sauce", "Garnishes"),
// each with their own list.
func FindIngredientLists(node *html.Node) []*html.Node {
	matcher := func(node *html.Node) (keep, exit bool) {
		if node.Type == html.ElementNode && node.DataAtom == atom.Ul {
			for _, a := range node.Attr {
				if a.Key == "class" && a.Val == "wprm-recipe-ingredients" {
					keep = true
					exit = true // No nested ingredients lists
				}
			}
		}
		return
	}
	return FindNodes(node, matcher)
}

// FindInstructionsList returns the node representing the list of recipe instructions.
// This implementation currently assumes that there is only 1 master list of instructions.
func FindInstructionsList(node *html.Node) *html.Node {
	return GetElementWithClass(node, atom.Ul, "wprm-recipe-instructions")
}

// GetElementWithClass returns the first element underneath and including `node` that has the given
// class value (as given in the HTML). The classes must be in the same order as those given.
func GetElementWithClass(node *html.Node, tagname atom.Atom, class string) *html.Node {
	matcher := func(node *html.Node) bool {
		if node.Type == html.ElementNode && node.DataAtom == tagname {
			for _, a := range node.Attr {
				if a.Key == "class" && a.Val == class {
					return true
				}
			}
		}
		return false
	}
	return FindNode(node, matcher)
}

// GetTextNode returns the first text node under the given node.
func GetTextNode(node *html.Node) *html.Node {
	matcher := func(node *html.Node) bool {
		return node.Type == html.TextNode
	}
	return FindNode(node, matcher)
}

// FindNode will return the first node that the matcher function accepts.
func FindNode(node *html.Node, matcher func(node *html.Node) bool) *html.Node {
	if matcher(node) {
		return node
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		n := FindNode(c, matcher)
		if n != nil {
			return n
		}
	}
	return nil
}

// FindNodes will gather all nodes that the matcher function accepts.
// The matcher function's second return value indicates whether to search the children of the
// current node.
// See TraverseNode from https://gist.github.com/Xeoncross/8bbb84bc4bf540bd907f79ee17c4e1fc
func FindNodes(node *html.Node, matcher func(node *html.Node) (bool, bool)) (nodes []*html.Node) {
	var keep, exit bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		keep, exit = matcher(n)
		if keep {
			nodes = append(nodes, n)
		}
		if exit {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	return nodes
}

func PrintNode(node *html.Node) {
	fmt.Print("Node Type: ")
	switch node.Type {
	case html.ElementNode:
		fmt.Println("Element")
	case html.TextNode:
		fmt.Println("Text")
	default:
		fmt.Println("Other")
	}

	fmt.Println("Node Data:", node.Data)

	fmt.Println("Node Attributes")
	for _, a := range node.Attr {
		fmt.Println(a.Key, a.Val)
	}
}
