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
func isTable(node *html.Node, class string) bool {
	return isElem(node, "tbody") && getAttr(node, "class") == class
}

func isPicture(node *html.Node, class string) bool {
	return isElem(node, "picture") && getAttr(node, "class") == class
}

func readItem(item *html.Node) *Item {
	time := item.FirstChild //.FirstChild
	temp := time.NextSibling.FirstChild.NextSibling.FirstChild

	return &Item{
		Ref:   getAttr(item, "span"),
		Title: "\n" + time.FirstChild.Data + " " + temp.Data + "°",
	}

}
func readItemMain(item *html.Node) *Item {
	temp := item.FirstChild //.FirstChild

	return &Item{
		Ref:   getAttr(item, "span"),
		Title: "Сейчас " + temp.Data + "°",
	}

}
func readItemWeekly(item *html.Node) *Item {
	date := item.FirstChild.FirstChild.FirstChild.NextSibling.FirstChild //.FirstChild
	temp := item.FirstChild.NextSibling.FirstChild.FirstChild.NextSibling.FirstChild
	log.Debug(date.Data)
	log.Debug(temp.Data)
	return &Item{
		Ref:   getAttr(item, "span"),
		Title: date.Data + " Сентября" + " " + temp.Data + "°",
	}

}

var a = 0
var items map[string]*Item
var items2 map[string]*Item

func searchFilmsNames(node *html.Node) map[string]*Item {
	if isDiv(node, "HhSR MBvM") {
		mainTemp := readItemMain(node)
		items[mainTemp.Title] = mainTemp

	}
	if isTable(node, "bT1T") {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if item := readItemWeekly(c); item != nil {
				items[item.Title] = item
			}
		}

	}

	if isDiv(node, "v8rM vk7t") {
		if item := readItem(node); item != nil {
			items[item.Title] = item

		}
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
	log.Info("sending request to https://weather.rambler.ru/v-moskve/")
	if response, err := http.Get("https://weather.rambler.ru/v-moskve/"); err != nil {
		log.Error("request to www.afusha.ru failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from https://weather.rambler.ru/v-moskve", "status", status)
		if status == http.StatusOK {
			if data, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from www.afisha.ru", "error", err)
			} else {
				log.Info("HTML from www.afisha.ru parsed successfully")
				items = make(map[string]*Item)
				items2 = make(map[string]*Item)
				searchFilmsNames(data)

				return items
			}
		}
	}
	return nil
}
