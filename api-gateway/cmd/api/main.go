package main

import (
	"api-gateway/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/routes"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	router := chi.NewRouter()
	APIHandlerAuth := handlers.NewAPIHandlerAuth(cfg)
	APIHandlerMovie := handlers.NewAPIHandlerMovie(cfg)

	route := routes.NewRoutes(router)

	route.SetupRouteAPIAuth(&APIHandlerAuth)
	route.SetupRouteAPIMovie(&APIHandlerMovie)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
