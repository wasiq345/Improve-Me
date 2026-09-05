package main

import (
	"context"
	"database/sql"
	"log"
	"myapp/internal/api"
	"myapp/internal/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	c := cron.New(cron.WithLocation(time.UTC))
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	const port = "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: api.EnableCORS(mux),
	}

	config := api.Config{
		Users: dbQueries,
		Notes: dbQueries,
	}

	c.AddFunc("@daily", func() {
		config.Users.ResetDailyCount(ctx)
	})
	c.Start()
	config.RegisterRoutes(mux)
	println("Server Running on port 8080")
	server.ListenAndServe()
}
