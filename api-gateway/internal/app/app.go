package app

import (
	"api-gateway/config"
	"api-gateway/config/cache"
	"api-gateway/internal/handlers"
	"api-gateway/internal/routes"
	"api-gateway/internal/server"
	"github.com/go-chi/chi/v5"
	"os"
)

func Run() {
	cfg := config.LoadConfig()
	redisClient := cache.NewRedisClient(os.Getenv("REDIS_HOST"), "", 0)

	router := chi.NewRouter()

	APIHandlerAuth := handlers.NewAPIHandlerAuth(cfg)
	APIHandlerMovie := handlers.NewAPIHandlerMovie(cfg, redisClient)
	APIHandlerUserList := handlers.NewAPIHandlerUserList(cfg)

	route := routes.NewRoutes(router)

	route.SetupRouteAPIAuth(&APIHandlerAuth)
	route.SetupRouteAPIMovie(&APIHandlerMovie)
	route.SetupRouteAPIUserList(&APIHandlerUserList)

	httpServer := server.NewServer(":8080", router)
	httpServer.Runserver()

}
