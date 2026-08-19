package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

var Links []Link

func main() {
	url := "http://127.0.0.1:5500/gophercises/htm-parser/index.html"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	processAllProducts(doc)

	fmt.Println(Links)
}

func processAllProducts(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		link := processNode(n)
		Links = append(Links, link)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processAllProducts(c)
	}
}

func processNode(n *html.Node) Link {
	for _, a := range n.Attr {
		if a.Key == "href" {
			link := Link{
				Href: a.Val,
				Text: n.FirstChild.Data,
			}
			return link
		}
	}

	return Link{}
}
