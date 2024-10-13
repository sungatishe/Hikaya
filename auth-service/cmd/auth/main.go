package main

import (
	"auth-service/config/db"
	"auth-service/internal/handlers"
	"auth-service/internal/repository"
	rabbitmq "auth-service/internal/rmq"
	"auth-service/internal/routes"
	"auth-service/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	db.InitDb()

	secretKey := os.Getenv("JWT_KEY")
	rabbit, err := rabbitmq.NewRabbitMQ(os.Getenv("RABBITMQ_URL"), "user_events")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	userRepo := repository.NewUserRepository(db.Db)
	userService := service.NewUserService(userRepo, secretKey, rabbit)
	userHandler := handlers.NewUserHandler(userService)

	router := chi.NewRouter()

	route := routes.NewRoutes(router)

	route.SetupRouteUser(userHandler)

	err = http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
