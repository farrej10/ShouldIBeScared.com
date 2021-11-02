package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	c := colly.NewCollector()

	c.OnHTML("tr > td > a", func(e *colly.HTMLElement) {
		fmt.Printf("Title Found: %q\n", e.Text)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("http://www.ifco.ie/en/IFCO/Pages/AllSearch/Start=1&Query=%5BApp_WorkTitle%5D=%22star%20wars%22&SearchTerm=star%20wars")
}

func main() {

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Println("Starting router")
	myRouter.HandleFunc("/", IndexHandler)

	http.ListenAndServe(":8089", myRouter)
}
