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
	address = "movie-service:50051"
)

type ViewData struct {
	Movie  Movie
	Movies []Movie
}

type ViewIndex struct {
	Popular  []Movie
	Trending []Movie
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

func GetMovie(c pb.MoviemangerClient, id string, rc chan Movie) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.GetMovie(ctx, &pb.Params{Id: id})
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

func GetRecommendations(c pb.MoviemangerClient, id string, rc chan []Movie) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resprec, err := c.GetRecommendations(ctx, &pb.Params{Id: id})
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

func GetPopular(c pb.MoviemangerClient, id string, rc chan []Movie) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resprec, err := c.GetPopular(ctx, &pb.Params{Id: id})
	if err != nil {
		log.Fatalln(err)
	}

	popular := []Movie{}
	for _, movie := range resprec.Movies {
		data := Movie{
			Title:       movie.Title,
			Description: movie.Description,
			Timestamp:   movie.Timestamp,
			ImageURL:    movie.Url,
			Id:          movie.Id,
		}
		popular = append(popular, data)
	}
	rc <- popular
}

func GetTrending(c pb.MoviemangerClient, id string, rc chan []Movie) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resprec, err := c.GetTrending(ctx, &pb.Params{Id: id})
	if err != nil {
		log.Fatalln(err)
	}

	trending := []Movie{}
	for _, movie := range resprec.Movies {
		data := Movie{
			Title:       movie.Title,
			Description: movie.Description,
			Timestamp:   movie.Timestamp,
			ImageURL:    movie.Url,
			Id:          movie.Id,
		}
		trending = append(trending, data)
	}
	rc <- trending
}

func (rw *RequestWrapper) IndexHandler(w http.ResponseWriter, r *http.Request) {

	popularChan := make(chan []Movie)
	trendingChan := make(chan []Movie)
	defer close(popularChan)
	defer close(trendingChan)

	go GetPopular(rw.c, mux.Vars(r)["id"], popularChan)
	go GetTrending(rw.c, mux.Vars(r)["id"], trendingChan)

	popular := <-popularChan
	trending := <-trendingChan

	err := rw.templates.Execute(w, ViewIndex{Popular: popular, Trending: trending})
	if err != nil {
		log.Fatalln(err)
	}
}

func (rw *RequestWrapper) MoviePageHandler(w http.ResponseWriter, r *http.Request) {

	movieChan := make(chan Movie)
	recommendChan := make(chan []Movie)
	defer close(movieChan)
	defer close(recommendChan)

	go GetMovie(rw.c, mux.Vars(r)["id"], movieChan)
	go GetRecommendations(rw.c, mux.Vars(r)["id"], recommendChan)

	movie := <-movieChan
	recommendations := <-recommendChan

	err := rw.templates.Execute(w, ViewData{Movie: movie, Movies: recommendations})
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

	templates := template.Must(template.ParseFiles("./pageGen/templates/movie.html", "./pageGen/templates/base.html"))
	templatesIndex := template.Must(template.ParseFiles("./pageGen/templates/index.html", "./pageGen/templates/base.html"))

	rw := &RequestWrapper{templates: templates, c: c}
	rwindex := &RequestWrapper{templates: templatesIndex, c: c}

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Println("Starting router")
	myRouter.HandleFunc("/", rwindex.IndexHandler)
	myRouter.HandleFunc("/movies/{id}", rw.MoviePageHandler)
	http.ListenAndServe(":8085", myRouter)
}
