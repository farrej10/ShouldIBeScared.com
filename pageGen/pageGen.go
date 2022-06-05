package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	pb "github.com/farrej10/ShouldIBeScared.com/movie"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

const (
	address = "movie-service:50051"
	scraper = "http://scraper-service:8089/"
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

type SearchResult struct {
	Results []Movie
	Pages   int32
}

type SearchView struct {
	Movies []Movie
	Query  string
	Pages  []int
	Page   int
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

func Sanitize(s *string) {

	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		log.Fatal(err)
	}
	*s = reg.ReplaceAllString(*s, "")
	*s = strings.Replace(*s, " ", "+", -1)
}

func Search(c pb.MoviemangerClient, query string, page string, rc chan SearchResult) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resprec, err := c.Search(ctx, &pb.Params{Id: query, Page: page})
	if err != nil {
		log.Fatalln(err)
	}

	search := []Movie{}
	for _, movie := range resprec.Results.Movies {
		data := Movie{
			Title:       movie.Title,
			Description: movie.Description,
			Timestamp:   movie.Timestamp,
			ImageURL:    movie.Url,
			Id:          movie.Id,
		}
		search = append(search, data)
	}
	rc <- SearchResult{Results: search, Pages: resprec.Pages}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (rw *RequestWrapper) SearchHandler(w http.ResponseWriter, r *http.Request) {
	searchString := r.URL.Query()["search"][0]
	Sanitize(&searchString)
	var page string
	if r.URL.Query()["page"] != nil {
		page = r.URL.Query()["page"][0]
	} else {
		page = "1"
	}
	searchChan := make(chan SearchResult)
	defer close(searchChan)

	go Search(rw.c, searchString, page, searchChan)

	result := <-searchChan

	tmp := min(int(5-result.Pages%5), int(result.Pages))
	pages := make([]int, tmp)
	for i := range pages {
		pages[i] = 1 + i
	}
	log.Println(searchString)
	p, _ := strconv.Atoi(page)
	searchString = strings.Replace(searchString, "+", " ", -1)
	err := rw.templates.Execute(w, SearchView{Movies: result.Results, Query: searchString, Pages: pages, Page: p - 1})
	if err != nil {
		log.Fatalln(err)
	}

}

func ScraperHandler(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("GET", scraper, nil)
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	w.WriteHeader(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to get scraper response")
	}
	fmt.Fprint(w, string(body))
}

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMoviemangerClient(conn)

	templates := template.Must(template.ParseFiles("/var/www/movie.html", "/var/www/base.html"))
	templatesIndex := template.Must(template.ParseFiles("/var/www/index.html", "/var/www/base.html"))
	templatesSearch := template.Must(template.ParseFiles("/var/www/search.html", "/var/www/base.html", "/var/www/searchitem.html"))

	rw := &RequestWrapper{templates: templates, c: c}
	rwindex := &RequestWrapper{templates: templatesIndex, c: c}
	rwsearch := &RequestWrapper{templates: templatesSearch, c: c}

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Println("Starting router")
	myRouter.HandleFunc("/", rwindex.IndexHandler)
	myRouter.HandleFunc("/movies/{id:[0-9]+}", rw.MoviePageHandler)
	myRouter.HandleFunc("/search", rwsearch.SearchHandler)
	myRouter.HandleFunc("/scraper", ScraperHandler)

	http.ListenAndServe(":8085", myRouter)
}
