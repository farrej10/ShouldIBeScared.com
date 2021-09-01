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

type Movie struct {
	Title       string
	Description string
	Timestamp   string
	ImageURL    string
}

type RequestWrapper struct {
	templates *template.Template
	id        string
	c         pb.MoviemangerClient
}

func (rw *RequestWrapper) RequestHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := rw.c.GetMovie(ctx, &pb.Params{Id: rw.id})
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

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMoviemangerClient(conn)

	templates := template.Must(template.ParseFiles("templates/movie.html", "templates/base.html"))
	rw := &RequestWrapper{templates: templates, id: "1", c: c}

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Println("Starting router")
	myRouter.HandleFunc("/", rw.RequestHandler)
	http.ListenAndServe(":8085", myRouter)
}
