package main

import (
	"strings"

	"golang.org/x/net/html"
)

func extractRelevantText(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var relevantNode *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if relevantNode != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "szakagazat" {
					relevantNode = n.Parent
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if relevantNode == nil {
		return "", nil
	}

	var sb strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				if a.Key == "class" && strings.HasPrefix(a.Val, "szakagazatnem") {
					return
				}
			}
		}
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(relevantNode)

	return strings.TrimSpace(sb.String()), nil
}
