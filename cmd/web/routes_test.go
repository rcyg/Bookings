package main

import (
	"Bookings/internal/config"
	"fmt"
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing; test passed
	default:
		t.Error(fmt.Sprintf("It not of type *chi.Mux, but is %T", v))
	}

}
