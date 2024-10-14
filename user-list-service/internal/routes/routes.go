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
	rt.r.Post("/userList", movieHandler.AddMovieToUserList)
	rt.r.Get("/userList/{id}", movieHandler.GetUsersMovieList)
	rt.r.Put("/userList/{userID}/{movieID}", movieHandler.UpdateMovieListType)
	rt.r.Delete("/userList/{userID}/{movieID}", movieHandler.DeleteMovieFromUserList)
}
