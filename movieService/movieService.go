package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	pb "github.com/farrej10/ShouldIBeScared.com/movie"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type MoviemangerServer struct {
	pb.UnimplementedMoviemangerServer
}

type Recommendations struct {
	Page    int `json:"page"`
	Results []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		GenreIds         []int   `json:"genre_ids"`
		ID               int     `json:"id"`
		MediaType        string  `json:"media_type"`
		Title            string  `json:"title"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		Popularity       float64 `json:"popularity"`
		PosterPath       string  `json:"poster_path"`
		ReleaseDate      string  `json:"release_date"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type Movie struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		PosterPath   string `json:"poster_path"`
		BackdropPath string `json:"backdrop_path"`
	} `json:"belongs_to_collection"`
	Budget int `json:"budget"`
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	ID                  int     `json:"id"`
	ImdbID              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	return os.Getenv(key)
}

func (s *MoviemangerServer) GetMovie(ctx context.Context, in *pb.Params) (*pb.Movie, error) {
	url := "https://api.themoviedb.org/3/movie/" + in.Id
	var bearer = "Bearer " + goDotEnvVariable("TOKEN")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var movie *Movie
	err = json.Unmarshal([]byte(body), &movie)
	if err != nil {
		log.Fatalln(err)
	}

	returnedmovie := &pb.Movie{
		Title:       movie.OriginalTitle,
		Description: movie.Overview,
		Timestamp:   movie.ReleaseDate,
		Url:         "https://image.tmdb.org/t/p/w1280" + movie.BackdropPath,
		Id:          strconv.Itoa(movie.ID),
	}

	return returnedmovie, nil
}

func (s *MoviemangerServer) GetRecommendations(ctx context.Context, in *pb.Params) (*pb.Movies, error) {

	var bearer = "Bearer " + goDotEnvVariable("TOKEN")
	client := &http.Client{Timeout: time.Second * 10}
	url := "https://api.themoviedb.org/3/movie/" + in.Id + "/recommendations"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	var recommendations *Recommendations
	err = json.Unmarshal([]byte(body), &recommendations)
	if err != nil {
		log.Fatalln(err)
	}

	movies := []*pb.Movie{}

	for _, element := range recommendations.Results {

		var tmpurl string

		if element.PosterPath == "" && element.BackdropPath == "" {
			log.Println("No Poster or Backdrop")
			tmpurl = ""
		} else if element.PosterPath == "" {
			log.Println("No Poster")
			tmpurl = "https://image.tmdb.org/t/p/w185" + element.BackdropPath
		} else {
			tmpurl = "https://image.tmdb.org/t/p/w154" + element.PosterPath
		}

		movie := pb.Movie{
			Title:       element.Title,
			Description: element.Overview,
			Timestamp:   element.ReleaseDate,
			Url:         tmpurl,
			Id:          strconv.Itoa(element.ID),
		}
		if movie.Url != "" {
			movies = append(movies, &movie)
		}

	}

	return &pb.Movies{Movies: movies}, nil
}

func (s *MoviemangerServer) GetPopular(ctx context.Context, in *pb.Params) (*pb.Movies, error) {

	var bearer = "Bearer " + goDotEnvVariable("TOKEN")
	client := &http.Client{Timeout: time.Second * 10}
	url := "https://api.themoviedb.org/3/movie/popular"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	var popular *Recommendations
	err = json.Unmarshal([]byte(body), &popular)
	if err != nil {
		log.Fatalln(err)
	}

	movies := []*pb.Movie{}

	for _, element := range popular.Results {

		var tmpurl string

		if element.PosterPath == "" && element.BackdropPath == "" {
			log.Println("No Poster or Backdrop")
			tmpurl = ""
		} else if element.PosterPath == "" {
			log.Println("No Poster")
			tmpurl = "https://image.tmdb.org/t/p/w185" + element.BackdropPath
		} else {
			tmpurl = "https://image.tmdb.org/t/p/w154" + element.PosterPath
		}

		movie := pb.Movie{
			Title:       element.Title,
			Description: element.Overview,
			Timestamp:   element.ReleaseDate,
			Url:         tmpurl,
			Id:          strconv.Itoa(element.ID),
		}
		if movie.Url != "" {
			movies = append(movies, &movie)
		}

	}

	return &pb.Movies{Movies: movies}, nil
}

func (s *MoviemangerServer) GetTrending(ctx context.Context, in *pb.Params) (*pb.Movies, error) {

	var bearer = "Bearer " + goDotEnvVariable("TOKEN")
	client := &http.Client{Timeout: time.Second * 10}
	url := "https://api.themoviedb.org/3/trending/movie/week"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	var trending *Recommendations
	err = json.Unmarshal([]byte(body), &trending)
	if err != nil {
		log.Fatalln(err)
	}

	movies := []*pb.Movie{}

	for _, element := range trending.Results {

		var tmpurl string

		if element.PosterPath == "" && element.BackdropPath == "" {
			log.Println("No Poster or Backdrop")
			tmpurl = ""
		} else if element.PosterPath == "" {
			log.Println("No Poster")
			tmpurl = "https://image.tmdb.org/t/p/w185" + element.BackdropPath
		} else {
			tmpurl = "https://image.tmdb.org/t/p/w154" + element.PosterPath
		}

		movie := pb.Movie{
			Title:       element.Title,
			Description: element.Overview,
			Timestamp:   element.ReleaseDate,
			Url:         tmpurl,
			Id:          strconv.Itoa(element.ID),
		}
		if movie.Url != "" {
			movies = append(movies, &movie)
		}

	}

	return &pb.Movies{Movies: movies}, nil
}

func main() {

	lis, _ := net.Listen("tcp", port)
	s := grpc.NewServer()
	log.Println("Starting Service")
	pb.RegisterMoviemangerServer(s, &MoviemangerServer{})
	log.Printf("Listen %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
