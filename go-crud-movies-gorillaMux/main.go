package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Name     string    `json:"name"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applications/json")
	var newMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&newMovie)
	newMovie.Id = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode(newMovie)
}

func deleteMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applications/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func updateMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applications/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var newMovie Movie
			x := json.NewDecoder(r.Body).Decode(&newMovie)
			fmt.Printf("r-Body %v and  x %v\n", &r.Body, x)
			newMovie.Id = params["id"]
			fmt.Printf("newMovie %v\n", newMovie.Director)
			fmt.Println(x)
			movies = append(movies, newMovie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

}

func main() {

	r := mux.NewRouter()

	movies = append(movies, Movie{Id: "1", Isbn: "42578", Name: "Movie1", Director: &Director{Firstname: "Jhon", Lastname: "Doe"}})
	movies = append(movies, Movie{Id: "2", Isbn: "37845", Name: "Movie2", Director: &Director{Firstname: "Da", Lastname: "Vinci"}})

	r.HandleFunc("/movies", getAllMovies).Methods("GET")
	r.HandleFunc("/movie/{id}", getMovieById).Methods("GET")
	r.HandleFunc("/create", createMovie).Methods("POST")
	r.HandleFunc("/update/{id}", updateMovieById).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteMovieById).Methods("DELETE")

	fmt.Printf("Server started on http://localhost:8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
