package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

func scrape() {
	c := colly.NewCollector(
		colly.AllowedDomains("imsdb.com"),
		colly.MaxDepth(3),
	)
	c.OnHTML("p > a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnHTML("tr > td > a[href]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), "/scripts") {
			e.Request.Visit(e.Attr("href"))
		}

	})
	c.OnHTML("pre", func(e *colly.HTMLElement) {
		splitURL := strings.Split(e.Request.URL.String(), "/")
		filename := splitURL[len(splitURL)-1]
		e.Response.Save("./scraperService/files/" + filename)
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://imsdb.com/all-scripts.html")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	go scrape()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>running scraper</h1>")
}

func main() {

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Println("Starting router")
	myRouter.HandleFunc("/", IndexHandler)

	http.ListenAndServe(":8089", myRouter)
}
