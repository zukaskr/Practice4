package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Movie struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Genre  string `json:"genre"`
	Budget int    `json:"budget"`
}

var db *sql.DB

func main() {
	dbHost := os.Getenv("DB_HOST") // Docker 服务名 "db"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable: ", err)
	}

	fmt.Println("Starting the Server...")

	r := mux.NewRouter()

	// 4. 定义路由
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, _ := db.Query("SELECT id, title, genre, budget FROM movies")
	var movies []Movie
	for rows.Next() {
		var m Movie
		rows.Scan(&m.ID, &m.Title, &m.Genre, &m.Budget)
		movies = append(movies, m)
	}
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var m Movie
	json.NewDecoder(r.Body).Decode(&m)
	_, err := db.Exec("INSERT INTO movies (title, genre, budget) VALUES ($1, $2, $3)", m.Title, m.Genre, m.Budget)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var m Movie
	json.NewDecoder(r.Body).Decode(&m)
	db.Exec("UPDATE movies SET title=$1, genre=$2, budget=$3 WHERE id=$4", m.Title, m.Genre, m.Budget, params["id"])
	w.WriteHeader(http.StatusOK)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db.Exec("DELETE FROM movies WHERE id=$1", params["id"])
	w.WriteHeader(http.StatusNoContent)
}
