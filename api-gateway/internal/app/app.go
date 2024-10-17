package app

import (
	"api-gateway/config"
	"api-gateway/config/cache"
	_ "api-gateway/docs"
	"api-gateway/internal/handlers"
	"api-gateway/internal/routes"
	"api-gateway/internal/server"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"
)

func Run() {
	cfg := config.LoadConfig()
	redisClient := cache.NewRedisClient(os.Getenv("REDIS_HOST"), "", 0)

	router := chi.NewRouter()

	APIHandlerAuth := handlers.NewAPIHandlerAuth(cfg)
	APIHandlerMovie := handlers.NewAPIHandlerMovie(cfg, redisClient)
	APIHandlerUserList := handlers.NewAPIHandlerUserList(cfg)
	APIHandlerRating := handlers.NewAPIHandlerRating(cfg)

	route := routes.NewRoutes(router)

	route.SetupRouteAPIAuth(&APIHandlerAuth)
	route.SetupRouteAPIMovie(&APIHandlerMovie)
	route.SetupRouteAPIUserList(&APIHandlerUserList)
	route.SetupRouteAPIRating(&APIHandlerRating)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	httpServer := server.NewServer(cfg.Port, router)
	httpServer.Runserver()

}
