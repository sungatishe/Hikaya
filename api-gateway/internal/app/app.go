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

	route := routes.NewRoutes(router)

	route.SetupRouteAPIAuth(&APIHandlerAuth)
	route.SetupRouteAPIMovie(&APIHandlerMovie)
	route.SetupRouteAPIUserList(&APIHandlerUserList)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	httpServer := server.NewServer(":8080", router)
	httpServer.Runserver()

}
