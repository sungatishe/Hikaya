package routes

import (
	"github.com/go-chi/chi/v5"
	"user-list-service/internal/handlers"
)

type Routes struct {
	r chi.Router
}

func NewRoutes(r chi.Router) *Routes {
	return &Routes{r}
}

func (rt *Routes) SetupRouteMovie(movieHandler *handlers.UserListHandler) {
	rt.r.Post("/add-to-list", movieHandler.AddMovieToUserList)
}
