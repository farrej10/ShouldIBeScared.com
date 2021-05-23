package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	PageTitle string
	MovieId   string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	data := Movie{
		PageTitle: "TEST TITLE",
	}

	t, _ := template.ParseFiles("templates/layout.html")
	t.Execute(w, data)
}

func movieHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]

	data := Movie{
		PageTitle: "TEST TITLE",
		MovieId:   key,
	}

	t, _ := template.ParseFiles("templates/layout.html")
	t.Execute(w, data)
}

func main() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", indexHandler)
	myRouter.HandleFunc("/movie/{id}", movieHandler)
	http.ListenAndServe(":8085", myRouter)
}
