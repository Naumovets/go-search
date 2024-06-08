package application

import (
	"github.com/Naumovets/go-search/internal/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(jsonMiddleware)
	router.Route("/api", a.loadSearchRoutes)

	return router
}

func (a *App) loadSearchRoutes(router chi.Router) {
	h := handler.New(a.rep)
	router.Get("/search", h.Search)
}
