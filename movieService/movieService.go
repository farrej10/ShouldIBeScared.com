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

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	port         = ":50051"
	cacheTimeout = 18000
)

type MoviemangerServer struct {
	pb.UnimplementedMoviemangerServer
	redisClient *redis.Client
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

func getVariable(key string) string {

	token := os.Getenv(key)
	if token == "" {
		log.Fatalf("error loading token")
	}

	return token
}

func SetMovieCache(client *redis.Client, body []byte, id string) {
	err := client.Set(id, []byte(body), cacheTimeout*time.Second).Err()
	if err != nil {
		log.Println(err)
	}
}

func (s *MoviemangerServer) GetMovie(ctx context.Context, in *pb.Params) (*pb.Movie, error) {

	redisJSON, err := s.redisClient.Get(in.Id).Result()
	if err != nil {

		url := "https://api.themoviedb.org/3/movie/" + in.Id
		var bearer = "Bearer " + getVariable("TOKEN")

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

		//add to redis before returning
		go SetMovieCache(s.redisClient, body, in.Id)

		return returnedmovie, nil
	}

	var movie *Movie
	err = json.Unmarshal([]byte(redisJSON), &movie)
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

func (recommendations *Recommendations) ParseMovieList() []*pb.Movie {
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
			tmpurl = "https://image.tmdb.org/t/p/w342" + element.PosterPath
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

	if len(movies) == 0 {
		url := "https://api.themoviedb.org/3/movie/popular"

		var popular *Recommendations
		_, err := GetMoviesHttp(url, &popular)
		if err != nil {
			log.Println("Error on response", err)
		}
		return popular.ParseMovieList()
	}

	return movies
}

func SetMovieListCache(client *redis.Client, body []byte, id string, listname string) {
	err := client.Set(id+listname, []byte(body), cacheTimeout*time.Second).Err()
	if err != nil {
		log.Println(err)
	}
	log.Println(listname + " Cache set")
}

func GetMoviesHttp(url string, r **Recommendations) ([]byte, error) {

	var bearer = "Bearer " + getVariable("TOKEN")
	client := &http.Client{Timeout: time.Second * 10}

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

	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		log.Fatalln(err)
	}

	return body, err
}

func (s *MoviemangerServer) GetRecommendations(ctx context.Context, in *pb.Params) (*pb.Movies, error) {

	redisJSON, err := s.redisClient.Get(in.Id + "recommend").Result()
	if err != nil {

		url := "https://api.themoviedb.org/3/movie/" + in.Id + "/recommendations"
		var recommendations *Recommendations
		body, err := GetMoviesHttp(url, &recommendations)
		if err != nil {
			log.Println("Error on response", err)
		}

		go SetMovieListCache(s.redisClient, body, in.Id, "recommend")

		return &pb.Movies{Movies: recommendations.ParseMovieList()}, nil
	}

	var recommendations *Recommendations
	err = json.Unmarshal([]byte(redisJSON), &recommendations)
	if err != nil {
		log.Fatalln(err)
	}

	return &pb.Movies{Movies: recommendations.ParseMovieList()}, nil
}

func (s *MoviemangerServer) GetPopular(ctx context.Context, in *pb.Params) (*pb.Movies, error) {
	redisJSON, err := s.redisClient.Get(in.Id + "popular").Result()
	if err != nil {

		url := "https://api.themoviedb.org/3/movie/popular"

		var popular *Recommendations
		body, err := GetMoviesHttp(url, &popular)
		if err != nil {
			log.Println("Error on response", err)
		}

		go SetMovieListCache(s.redisClient, body, in.Id, "popular")

		return &pb.Movies{Movies: popular.ParseMovieList()}, nil
	}

	var popular *Recommendations
	err = json.Unmarshal([]byte(redisJSON), &popular)
	if err != nil {
		log.Fatalln(err)
	}

	return &pb.Movies{Movies: popular.ParseMovieList()}, nil
}

func (s *MoviemangerServer) GetTrending(ctx context.Context, in *pb.Params) (*pb.Movies, error) {

	redisJSON, err := s.redisClient.Get(in.Id + "trending").Result()
	if err != nil {

		url := "https://api.themoviedb.org/3/trending/movie/week"
		var trending *Recommendations
		body, err := GetMoviesHttp(url, &trending)
		if err != nil {
			return nil, err
		}
		go SetMovieListCache(s.redisClient, body, in.Id, "trending")
		return &pb.Movies{Movies: trending.ParseMovieList()}, nil
	}
	var trending *Recommendations
	err = json.Unmarshal([]byte(redisJSON), &trending)
	if err != nil {
		log.Println("Error on response", err)
	}

	return &pb.Movies{Movies: trending.ParseMovieList()}, nil
}

func (s *MoviemangerServer) Search(ctx context.Context, in *pb.Params) (*pb.SearchResult, error) {

	redisJSON, err := s.redisClient.Get(in.Id + in.Page + "search").Result()
	if err != nil {

		url := "https://api.themoviedb.org/3/search/movie?language=en-US&page=" + in.Page + "&include_adult=false&query=" + in.Id

		var searchResults *Recommendations
		body, err := GetMoviesHttp(url, &searchResults)
		if err != nil {
			return nil, err
		}
		go SetMovieListCache(s.redisClient, body, in.Id, in.Page+"search")
		return &pb.SearchResult{Results: &pb.Movies{Movies: searchResults.ParseMovieList()}, Pages: int32(searchResults.TotalPages)}, nil
	}
	var searchResults *Recommendations
	err = json.Unmarshal([]byte(redisJSON), &searchResults)
	if err != nil {
		log.Println("Error on response", err)
	}

	return &pb.SearchResult{Results: &pb.Movies{Movies: searchResults.ParseMovieList()}, Pages: int32(searchResults.TotalPages)}, nil
}

func main() {

	lis, _ := net.Listen("tcp", port)
	s := grpc.NewServer()
	log.Println("Starting Service")

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	log.Println("Connecting to Redis")
	pb.RegisterMoviemangerServer(s, &MoviemangerServer{redisClient: client})
	log.Printf("Listen %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
