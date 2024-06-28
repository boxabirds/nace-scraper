package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type NACECode struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Code          string     `json:"code"`
	Description   string     `json:"description"`
	Href          string     `json:"href"`
	Title         string     `json:"title"`
	Level         int        `json:"level"`
	SubCategories []Category `json:"sub_categories,omitempty"`
}

func fetchAndParseNACECodes(url string) (NACECode, error) {
	resp, err := http.Get(url)
	if err != nil {
		return NACECode{}, fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NACECode{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return NACECode{}, fmt.Errorf("failed to parse HTML: %v", err)
	}

	naceCode := NACECode{}

	// Find the div with class "nacelist"
	var nacelistDiv *html.Node
	var findNacelist func(*html.Node)
	findNacelist = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "nacelist") {
			nacelistDiv = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNacelist(c)
		}
	}
	findNacelist(doc)

	if nacelistDiv == nil {
		return NACECode{}, fmt.Errorf("could not find div with class 'nacelist'")
	}

	naceCode.Categories = parseCategories(nacelistDiv, 1)

	return naceCode, nil
}

func parseCategories(n *html.Node, level int) []Category {
	var categories []Category

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "ul" && hasClass(node, fmt.Sprintf("level%d", level)) {
			for li := node.FirstChild; li != nil; li = li.NextSibling {
				if li.Type == html.ElementNode && li.Data == "li" && hasClass(li, fmt.Sprintf("level%d", level)) {
					category := parseCategory(li, level)
					categories = append(categories, category)
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return categories
}

func parseCategory(li *html.Node, level int) Category {
	category := Category{Level: level}

	for node := li.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.ElementNode && node.Data == "a" {
			category.Href = getAttr(node, "href")
			category.Title = getAttr(node, "title")
			if node.FirstChild != nil {
				category.Code, category.Description = parseCodeAndDescription(node.FirstChild.Data)
			}
		} else if node.Type == html.ElementNode && node.Data == "ul" && hasClass(node, fmt.Sprintf("level%d", level+1)) {
			category.SubCategories = parseCategories(node, level+1)
		}
	}

	return category
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, c := range classes {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func parseCodeAndDescription(s string) (string, string) {
	parts := strings.SplitN(s, " - ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return "", s
}

func printCategories(categories []Category, indent int) {
	for _, cat := range categories {
		fmt.Printf("%s%s - %s (Level: %d)\n", strings.Repeat("  ", indent), cat.Code, cat.Description, cat.Level)
		fmt.Printf("%sHref: %s\n", strings.Repeat("  ", indent+1), cat.Href)
		fmt.Printf("%sTitle: %s\n", strings.Repeat("  ", indent+1), cat.Title)
		if len(cat.SubCategories) > 0 {
			printCategories(cat.SubCategories, indent+1)
		}
	}
}

func main() {
	url := "https://companyformationhungary.com/nace-codes.html"
	naceCode, err := fetchAndParseNACECodes(url)
	if err != nil {
		log.Fatalf("Error fetching and parsing NACE codes: %v", err)
	}

	printCategories(naceCode.Categories, 0)
}
