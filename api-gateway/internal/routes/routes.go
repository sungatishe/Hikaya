package routes

import (
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Routes struct {
	r chi.Router
}

func NewRoutes(r chi.Router) *Routes {
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
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
	rt.r.Post("/logout", APIHandlerAuth.LogOut)
	rt.r.Get("/user", (APIHandlerAuth.User))
}

func (rt *Routes) SetupRouteAPIUserList(APIHandlerUserList *handlers.APIHandlerUserList) {
	rt.r.Post("/userList", middleware.AuthMiddleware(APIHandlerUserList.AddMovieToUsersList))
	rt.r.Get("/userList/{id}", middleware.AuthMiddleware(APIHandlerUserList.GetUsersList))
	rt.r.Put("/userList/{userID}/{movieID}", middleware.AuthMiddleware(APIHandlerUserList.UpdateMovieListType))
	rt.r.Delete("/userList/{userID}/{movieID}", middleware.AuthMiddleware(APIHandlerUserList.DeleteMovieFromUserList))
}

func (rt *Routes) SetupRouteAPIRating(APIHandlerRating *handlers.APIHandlerRating) {
	rt.r.Post("/reviews", middleware.AuthMiddleware(APIHandlerRating.CreateMovieReview))
	rt.r.Get("/movie/{id}/reviews", APIHandlerRating.GetMovieReviews)
	rt.r.Get("/movie/{id}/rating", APIHandlerRating.GetMovieRating)
}
