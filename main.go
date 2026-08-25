package main

import (
	"database/sql"
	"log"
	"myapp/internal/api"
	"myapp/internal/database"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	const port = "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	config := api.Config{
		DB: dbQueries,
	}

	config.RegisterRoutes(mux)
	println("Server Running on port 8080")
	server.ListenAndServe()
}
