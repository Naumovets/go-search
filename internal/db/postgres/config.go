package postgres

import (
	"fmt"
	"os"

	"github.com/Naumovets/go-search/internal/db"
	"github.com/joho/godotenv"
)

func NewConfig(path string) (*db.Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("cannot load env from %s", path)
	}

	return &db.Config{
		DB_HOST:  os.Getenv("DB_HOST"),
		DB_PORT:  os.Getenv("DB_PORT"),
		DB_NAME:  os.Getenv("DB_NAME"),
		PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		USER:     os.Getenv("POSTGRES_USER"),
	}, nil

}
