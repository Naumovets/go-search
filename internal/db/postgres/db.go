package postgres

import (
	"fmt"

	"github.com/Naumovets/go-search/internal/db"
	"github.com/jmoiron/sqlx"
)

func NewConn(cfg db.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.USER, cfg.PASSWORD, cfg.DB_NAME, cfg.DB_PORT)

	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("sqlx err: %s", err)
	}

	err = conn.Ping()

	if err != nil {
		return nil, fmt.Errorf("ping failed: %s", err)
	}

	return conn, nil
}
