package api

import "myapp/internal/database"

type Config struct {
	DB *database.Queries
}
