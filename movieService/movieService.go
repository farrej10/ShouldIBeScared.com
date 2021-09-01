package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/farrej10/ShouldIBeScared.com/movie"
	"github.com/joho/godotenv"
	"github.com/valyala/fastjson"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type MovieData struct {
	Title       string
	Description string
	Timestamp   string
	ImageURL    string
}

type MoviemangerServer struct {
	pb.UnimplementedMoviemangerServer
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

	// log.Println(string([]byte(body)))

	var p fastjson.Parser
	v, err := p.Parse(string([]byte(body)))
	if err != nil {
		log.Fatal(err)
	}

	data := MovieData{
		Title:       string(v.GetStringBytes("original_title")),
		Description: string(v.GetStringBytes("overview")),
		Timestamp:   "updated now!",
		ImageURL:    "https://image.tmdb.org/t/p/w1280" + string(v.GetStringBytes("backdrop_path")),
	}
	log.Println("Sending Data Now")
	return &pb.Movie{Title: data.Title, Description: data.Description, Timestamp: data.Timestamp, Url: data.ImageURL}, nil
}

func main() {

	lis, _ := net.Listen("tcp", port)
	s := grpc.NewServer()

	pb.RegisterMoviemangerServer(s, &MoviemangerServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
