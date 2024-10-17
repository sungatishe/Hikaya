package service

import (
	"rating-service/internal/models"
	"rating-service/internal/repository"
)

type ratingService struct {
	repo repository.RatingRepository
}

func NewRatingService(repo repository.RatingRepository) RatingService {
	return &ratingService{repo: repo}
}

func (r ratingService) CreateMovieReview(review *models.Review) error {
	return r.repo.CreateMovieReview(review)
}

func (r ratingService) GetReviewsByMovieID(movieID uint) ([]models.Review, error) {
	return r.repo.GetReviewsByMovieID(movieID)
}

func (r ratingService) CalculateMovieRating(movieID uint) (float64, error) {
	return r.repo.CalculateMovieRating(movieID)
}

func (r ratingService) UpdateMovieRating(movieID uint, avrRating float64, reviewCount int) error {
	return r.repo.UpdateMovieRating(movieID, avrRating, reviewCount)
}
