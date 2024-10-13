package main

import (
	"github.com/go-chi/chi/v5"
	"movie-service/config/db"
	"movie-service/internal/handlers"
	"movie-service/internal/rabbitMQ"
	"movie-service/internal/repository"
	"movie-service/internal/routes"
	"movie-service/internal/service"
	"net/http"
)

func main() {
	db.InitDb()

	movieRepo := repository.NewMovieRepository(db.Db)
	movieService := service.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)
	rmqConsumer := rabbitMQ.NewConsumer("amqp://guest:guest@rabbitmq:5672/")

	// Запуск консюмера в отдельной горутине
	go rmqConsumer.Consume("movie_delete", movieHandler.HandleDeleteMovieEvent)
	go rmqConsumer.Consume("movie_create", movieHandler.HandleCreateMovieEvent)

	router := chi.NewRouter()

	route := routes.NewRoutes(router)

	route.SetupRouteMovie(movieHandler)

	err := http.ListenAndServe(":8082", router)
	if err != nil {
		panic(err)
	}
}
