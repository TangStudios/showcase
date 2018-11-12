package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"golang.org/x/net/html"
)

type Order struct {
	Size int
	List []int
}

type Showcase struct {
	Title  string
	Studio string
	Link   string
}

var URL = "https://unity.com"

var (
	showcases  []Showcase
	userOrders map[string]*Order
)

func main() {
	userOrders = make(map[string]*Order)

	resp, getErr := http.Get(URL + "/madewith")
	if getErr != nil {
		panic("error getting page info")
	}
	root, parseErr := html.Parse(resp.Body)
	if parseErr != nil {
		panic("error parsing response")
	}
	scrapePage(root)

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":80", nil)
}

func scrapePage(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "article" {
		// hard-coded web scraping from https://unity.com/madewith
		// scrapes out the title, studio, and link to more info of a game
		dataNode := n.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling
		var showcase Showcase
		showcase.Title = dataNode.FirstChild.NextSibling.FirstChild.Data
		showcase.Studio = dataNode.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.Data
		showcase.Link = URL + n.FirstChild.NextSibling.Attr[0].Val

		fmt.Println(showcase.Title + " " + showcase.Studio + " " + showcase.Link)

		showcases = append(showcases, showcase)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrapePage(c)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	order, ok := userOrders[r.RemoteAddr]
	if !ok {
		var newOrder Order
		var list []int

		for i := 0; i < len(showcases); i++ {
			list = append(list, i)
		}
		newOrder.List = list
		newOrder.Size = len(showcases)
		userOrders[r.RemoteAddr] = &newOrder
	} else if order.Size == 0 {
		order.Size = len(order.List)
	}

	userOrder := userOrders[r.RemoteAddr].List
	size := userOrders[r.RemoteAddr].Size
	rand := rand.Intn(size)
	showcase := showcases[userOrder[rand]]

	fmt.Printf("rand: %d  val: %d  size: %d\n", rand, userOrder[rand], size)
	fmt.Fprintf(w, "<p>%s</p><p>%s</p><a href=\"%s\">%s</a>", showcase.Title, showcase.Studio, showcase.Link, showcase.Link)

	// swap index to end
	temp := userOrder[rand]
	userOrder[rand] = userOrder[size-1]
	userOrder[size-1] = temp
	userOrders[r.RemoteAddr].Size--
	fmt.Println(userOrder)
}
