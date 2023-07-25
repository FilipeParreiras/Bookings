package main

import (
	"net/http"

	"github.com/FilipeParreiras/Bookings/pkg/config"
	"github.com/FilipeParreiras/Bookings/pkg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Recoverer) // Dont let app panic
	mux.Use(WriteToConsole)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
