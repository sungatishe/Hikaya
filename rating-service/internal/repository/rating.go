package repository

import (
	"gorm.io/gorm"
	"math"
	"rating-service/internal/models"
)

type ratingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{db}
}

func (r *ratingRepository) CreateMovieReview(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ratingRepository) GetReviewsByMovieID(movieID uint) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Where("movie_id = ?", movieID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ratingRepository) CalculateMovieRating(movieID uint) (float64, error) {
	var reviews []models.Review
	if err := r.db.Where("movie_id = ?", movieID).Find(&reviews).Error; err != nil {
		return 0, err
	}

	if len(reviews) == 0 {
		return 0, nil
	}

	totalRating := 0
	for _, review := range reviews {
		totalRating += review.Rating
	}

	avgRating := float64(totalRating) / float64(len(reviews))
	roundedRating := math.Round(avgRating*100) / 100
	return roundedRating, nil
}

func (r *ratingRepository) UpdateMovieRating(movieID uint, avrRating float64, reviewCount int) error {
	var movieRating models.MovieRating

	err := r.db.Where("movie_id = ?", movieID).First(&movieRating).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			movieRating = models.MovieRating{
				MovieID:     movieID,
				AvrRating:   avrRating,
				ReviewCount: reviewCount,
			}
		}
		return err
	}
	movieRating.AvrRating = avrRating
	movieRating.ReviewCount = reviewCount
	return r.db.Save(&movieRating).Error
}
