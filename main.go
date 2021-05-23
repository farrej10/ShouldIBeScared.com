package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	PageTitle string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	data := Movie{
		PageTitle: "TEST TITLE",
	}

	t, _ := template.ParseFiles("templates/layout.html")
	t.Execute(w, data)
}

func main() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8085", myRouter)
}
