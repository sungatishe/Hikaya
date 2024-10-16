package service

import "rating-service/internal/models"

type RatingService interface {
	CreateMovieReview(review *models.Review) error
	GetReviewsByMovieID(movieID uint) ([]models.Review, error)
	CalculateMovieRating(movieID uint) (float64, error)
	UpdateMovieRating(movieID uint, avrRating float64, reviewCount int) error
}
