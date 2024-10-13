package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"user-list-service/config/db"
	"user-list-service/internal/handlers"
	"user-list-service/internal/repository"
	"user-list-service/internal/routes"
	"user-list-service/internal/service"
)

func main() {
	db.InitDb()

	userListRepo := repository.NewUserListRepository(db.Db)
	userListService := service.NewUserService(userListRepo)
	userListHandler := handlers.NewUserListHandler(userListService)

	router := chi.NewRouter()

	route := routes.NewRoutes(router)

	route.SetupRouteMovie(userListHandler)

	err := http.ListenAndServe(":8083", router)
	if err != nil {
		panic(err)
	}
}
