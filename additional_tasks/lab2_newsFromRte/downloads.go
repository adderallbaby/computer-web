package main

import (
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"net/http"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}
func isSpan(node *html.Node, class string) bool {
	return isElem(node, "span") && getAttr(node, "class") == class
}

func isPicture(node *html.Node, class string) bool {
	return isElem(node, "img") && getAttr(node, "alt") == class
}

func readItem(item *html.Node) *Item {
	log.Debug(item.FirstChild.Data)

	//log.Debug(item.FirstChild.FirstChild.FirstChild.FirstChild.Data)
	text := item.FirstChild
	//img := item.Parent.Parent.Parent.PrevSibling.FirstChild.FirstChild.NextSibling
	//imageSrc = getAttr(item.Parent.Parent.Parent.Parent.PrevSibling.FirstChild, "href")
	//log.Info(a.Data)
	//log.Info(b.Data)

	return &Item{
		Ref:   getAttr(text, "href"),
		Title: text.Data,
	}
}

var items map[string]*Item
var imageSrc string

func searchFilmsNames(node *html.Node) map[string]*Item {

	if isSpan(node, "underline") {

		if item := readItem(node); item != nil {
			items[item.Title] = item

		}
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items1 := searchFilmsNames(c); items1 != nil && c.NextSibling == nil {
			return items1
		}
	}
	return nil
}

type Item struct {
	Ref, Title, ImageSrc string
}

func downloadNews() map[string]*Item {
	log.Info("sending request to www.afisha.ru")
	if response, err := http.Get("https://www.rte.ie/news/"); err != nil {
		log.Error("request to www.afusha.ru failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from www.afisha.ru", "status", status)
		if status == http.StatusOK {
			if data, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from www.afisha.ru", "error", err)
			} else {
				log.Info("HTML from www.afisha.ru parsed successfully")
				items = make(map[string]*Item)
				searchFilmsNames(data)
				return items
			}
		}
	}
	return nil
}
