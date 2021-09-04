package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	pb "github.com/farrej10/ShouldIBeScared.com/movie"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

type ViewData struct {
	Movie  Movie
	Movies []Movie
}

type Movie struct {
	Title       string
	Description string
	Timestamp   string
	ImageURL    string
	Id          string
}

type RequestWrapper struct {
	templates *template.Template
	c         pb.MoviemangerClient
}

func GetMovie(rw *RequestWrapper, id string, rc chan Movie) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := rw.c.GetMovie(ctx, &pb.Params{Id: id})
	if err != nil {
		log.Fatalln(err)
	}

	movie := Movie{
		Title:       resp.Title,
		Description: resp.Description,
		Timestamp:   resp.Timestamp,
		ImageURL:    resp.Url,
		Id:          resp.Id,
	}

	rc <- movie
}

func GetRecommendations(rw *RequestWrapper, id string, rc chan []Movie) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resprec, err := rw.c.GetRecommendations(ctx, &pb.Params{Id: id})
	if err != nil {
		log.Fatalln(err)
	}

	recommendations := []Movie{}
	for _, movie := range resprec.Movies {
		data := Movie{
			Title:       movie.Title,
			Description: movie.Description,
			Timestamp:   movie.Timestamp,
			ImageURL:    movie.Url,
			Id:          movie.Id,
		}
		recommendations = append(recommendations, data)
	}
	rc <- recommendations
}
func (rw *RequestWrapper) IndexHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := rw.c.GetMovie(ctx, &pb.Params{Id: mux.Vars(r)["id"]})
	if err != nil {
		log.Fatalln(err)
	}
	data := Movie{
		Title:       resp.Title,
		Description: resp.Description,
		Timestamp:   resp.Timestamp,
		ImageURL:    resp.Url,
	}
	err = rw.templates.Execute(w, data)
	if err != nil {
		log.Fatalln(err)
	}
}

func (rw *RequestWrapper) MoviePageHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	movieChan := make(chan Movie)
	recommendChan := make(chan []Movie)
	defer close(movieChan)
	defer close(recommendChan)

	go GetMovie(rw, mux.Vars(r)["id"], movieChan)
	go GetRecommendations(rw, mux.Vars(r)["id"], recommendChan)

	movie := <-movieChan
	recommendations := <-recommendChan

	err := rw.templates.Execute(w, ViewData{Movie: movie, Movies: recommendations})
	end := time.Now()
	log.Printf("Order processed after %v seconds\n", end.Sub(start).Seconds())
	log.Printf("Sending page: %v\n", movie.Id)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMoviemangerClient(conn)

	templates := template.Must(template.ParseFiles("templates/movie.html", "templates/base.html"))
	rw := &RequestWrapper{templates: templates, c: c}
	r2 := &RequestWrapper{templates: templates, c: c}

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Println("Starting router")
	myRouter.HandleFunc("/", rw.IndexHandler)
	myRouter.HandleFunc("/movies/{id}", r2.MoviePageHandler)
	http.ListenAndServe(":8085", myRouter)
}
