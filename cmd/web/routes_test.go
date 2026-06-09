package main

import (
	"testing"

	"github.com/dh-n/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("The return type should be *chi.Mux, instead received %T", v)
	}
}
