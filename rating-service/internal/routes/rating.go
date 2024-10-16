package routes

import (
	"github.com/go-chi/chi/v5"
	"rating-service/internal/handlers"
)

type Routes struct {
	r chi.Router
}

func NewRoutes(r chi.Router) *Routes {
	return &Routes{r}
}

func (rt *Routes) SetupRatingRoutes(ratingHandler *handlers.RatingHandler) {
	rt.r.Post("/reviews", ratingHandler.CreateMovieReview)
	rt.r.Get("/movie/{id}/reviews", ratingHandler.GetReviewsByMovieID)
	rt.r.Get("/movie/{id}/rating", ratingHandler.GetMovieRating)
}
