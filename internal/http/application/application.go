package application

import (
	"fmt"
	"net/http"

	"github.com/Naumovets/go-search/internal/repositories/db"
)

type App struct {
	router http.Handler
	rep    *db.Repository
}

func New(rep *db.Repository) *App {
	app := &App{
		rep: rep,
	}

	app.router = app.loadRoutes()

	return app
}

func (a *App) Start() error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := server.ListenAndServe()

	if err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
