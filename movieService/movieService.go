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

	pb "github.com/farrej10/ShouldIBeScared.com/movie"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Movie struct {
	Title       string
	Description string
	Timestamp   string
	ImageURL    string
	Id          string
}

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

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	return os.Getenv(key)
}

func (s *MoviemangerServer) GetMovie(ctx context.Context, in *pb.Params) (*pb.Movie, error) {
	url := "https://api.themoviedb.org/3/movie/268"
	var bearer = "Bearer " + goDotEnvVariable("TOKEN")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var reading map[string]interface{}
	err = json.Unmarshal([]byte(body), &reading)
	if err != nil {
		log.Fatalln(err)
	}

	data := Movie{
		Title:       reading["original_title"].(string),
		Description: reading["overview"].(string),
		Timestamp:   "updated now!",
		ImageURL:    "https://image.tmdb.org/t/p/w1280" + reading["backdrop_path"].(string),
		Id:          strconv.Itoa(int(reading["id"].(float64))),
	}

	log.Println("Sending Data Now")

	return &pb.Movie{Title: data.Title, Description: data.Description, Timestamp: data.Timestamp, Url: data.ImageURL}, nil
}

func (s *MoviemangerServer) GetMovies(ctx context.Context, in *pb.Params) (*pb.Movies, error) {
	url := "https://api.themoviedb.org/3/movie/" + in.Id
	var bearer = "Bearer " + goDotEnvVariable("TOKEN")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var reading map[string]interface{}
	err = json.Unmarshal([]byte(body), &reading)
	if err != nil {
		log.Fatalln(err)
	}

	data := pb.Movie{
		Title:       reading["original_title"].(string),
		Description: reading["overview"].(string),
		Timestamp:   "updated now!",
		Url:         "https://image.tmdb.org/t/p/w1280" + reading["backdrop_path"].(string),
		Id:          strconv.Itoa(int(reading["id"].(float64))),
	}

	url = "https://api.themoviedb.org/3/movie/" + in.Id + "/recommendations"

	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	// Send req using http Client
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	var recommendations *Recommendations
	err = json.Unmarshal([]byte(body), &recommendations)
	if err != nil {
		log.Fatalln(err)
	}

	movies := []*pb.Movie{&data}

	// for _, element := range recommendations["results"].([]interface{}) {
	// 	movie := pb.Movie{
	// 		Title:       element.(map[string]interface{})["title"].(string),
	// 		Description: element.(map[string]interface{})["overview"].(string),
	// 		Timestamp:   "updated now!",
	// 		Url:         "https://image.tmdb.org/t/p/w342" + element.(map[string]interface{})["poster_path"].(string),
	// 		Id:          strconv.Itoa(int(element.(map[string]interface{})["id"].(float64))),
	// 	}
	// 	movies = append(movies, &movie)
	// }

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
			Id:          strconv.Itoa(int(element.ID)),
		}
		if movie.Url != "" {
			movies = append(movies, &movie)
		}

	}
	log.Println("Sending Data Now")

	return &pb.Movies{Movies: movies}, nil
}

func main() {

	lis, _ := net.Listen("tcp", port)
	s := grpc.NewServer()
	log.Println("Starting Service")
	pb.RegisterMoviemangerServer(s, &MoviemangerServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
