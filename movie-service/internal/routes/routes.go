package routes

import (
	"github.com/go-chi/chi/v5"
	"movie-service/internal/handlers"
)

type Routes struct {
	r chi.Router
}

func NewRoutes(r chi.Router) *Routes {
	return &Routes{r}
}

func (rt *Routes) SetupRouteMovie(movieHandler *handlers.MovieHandler) {
	rt.r.Post("/movies", movieHandler.CreateMovie)
	rt.r.Get("/movies/{id}", movieHandler.GetMovieById)
	rt.r.Get("/movies", movieHandler.GetAllMovies)
	rt.r.Get("/movies/search", movieHandler.SearchMovies)
	rt.r.Put("/movies/{id}", movieHandler.UpdateMovie)
	rt.r.Delete("/movies/{id}", movieHandler.DeleteMovie)
}
