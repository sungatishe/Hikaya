package routes

import (
	"api-gateway/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Routes struct {
	r chi.Router
}

func NewRoutes(r chi.Router) *Routes {
	return &Routes{r}
}

func (rt *Routes) SetupRouteAPIMovie(APIHandlerMovie *handlers.APIHandlerMovie) {
	rt.r.Post("/movies", APIHandlerMovie.CreateMovieMQ)
	rt.r.Get("/movies", APIHandlerMovie.GetMovies)
	rt.r.Get("/movies/{id}", APIHandlerMovie.GetMovieByID)
	rt.r.Put("/movies/{id}", APIHandlerMovie.UpdateMovie)
	rt.r.Delete("/movies/{id}", APIHandlerMovie.DeleteMovieMQ)
}

func (rt *Routes) SetupRouteAPIAuth(APIHandlerAuth *handlers.APIHandlerAuth) {
	rt.r.Post("/register", APIHandlerAuth.RegisterUser)
	rt.r.Post("/login", APIHandlerAuth.LoginUser)
}
