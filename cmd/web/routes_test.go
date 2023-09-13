package main

import (
	"fmt"
	"github.com/FilipeParreiras/Bookings/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	// Do nothing
	default:
		t.Error(fmt.Sprintf("type is not chi.Mux type is instaad %T", v))
	}
}
