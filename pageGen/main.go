package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	Title       string
	Description string
	Timestamp   string
	ImageURL    string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	data := Movie{
		Title:       "TEST TITLES",
		Description: "testing this thing saying something else",
		Timestamp:   "updated now!",
		ImageURL:    "https://image.tmdb.org/t/p/w1280/nvjcKJCDPU9bDEEyTyneTj4PnuO.jpg",
	}

	t, err := template.ParseFiles("templates/movie.html", "templates/base.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

// func movieHandler(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	key := vars["id"]

// 	data := Movie{
// 		PageTitle: "TEST TITLE",
// 		MovieId:   key,
// 	}

// 	t, _ := template.ParseFiles("templates/layout.html")
// 	t.Execute(w, data)
// }

func main() {

	myRouter := mux.NewRouter().StrictSlash(true)
	fmt.Println("Starting router")
	myRouter.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8085", myRouter)
}
