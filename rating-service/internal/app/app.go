package app

import (
	"github.com/go-chi/chi/v5"
	"rating-service/config/db"
	"rating-service/internal/handlers"
	"rating-service/internal/repository"
	"rating-service/internal/routes"
	"rating-service/internal/server"
	"rating-service/internal/service"
)

func Run() {
	db.InitDB()

	ratingRepo := repository.NewRatingRepository(db.DB)
	ratingService := service.NewRatingService(ratingRepo)
	ratingHandler := handlers.NewRatingHandler(ratingService)

	router := chi.NewRouter()

	route := routes.NewRoutes(router)
	route.SetupRatingRoutes(ratingHandler)

	httpServer := server.NewServer("8084", router)
	httpServer.RunServer()
}
