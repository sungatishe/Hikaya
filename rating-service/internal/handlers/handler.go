package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"rating-service/internal/models"
	"rating-service/internal/service"
	"strconv"
)

type RatingHandler struct {
	service service.RatingService
}

func NewRatingHandler(service service.RatingService) *RatingHandler {
	return &RatingHandler{service: service}
}

func (rh *RatingHandler) CreateMovieReview(rw http.ResponseWriter, r *http.Request) {
	var review models.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	if review.Rating < 1 || review.Rating > 10 {
		http.Error(rw, "Rating must be between 1 and 10", http.StatusBadRequest)
		return
	}

	err = rh.service.CreateMovieReview(&review)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	avrRating, err := rh.service.CalculateMovieRating(review.MovieID)
	if err != nil {
		http.Error(rw, "Server error", http.StatusInternalServerError)
		return
	}

	reviews, err := rh.service.GetReviewsByMovieID(review.MovieID)
	if err != nil {
		http.Error(rw, "Server error", http.StatusInternalServerError)
		return
	}
	reviewCount := len(reviews)
	err = rh.service.UpdateMovieRating(review.MovieID, avrRating, reviewCount)
	if err != nil {
		http.Error(rw, "Server error", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(review)
}

func (rh *RatingHandler) GetReviewsByMovieID(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	reviews, err := rh.service.GetReviewsByMovieID(uint(id))
	if err != nil {
		http.Error(rw, "Reviews not found", http.StatusBadRequest)
		return
	}

	if len(reviews) == 0 {
		http.Error(rw, "Movie not found", http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(reviews)
}

func (rh *RatingHandler) GetMovieRating(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	avrRating, err := rh.service.CalculateMovieRating(uint(id))
	if err != nil {
		http.Error(rw, "Server err", http.StatusInternalServerError)
		return
	}

	if avrRating == 0 {
		http.Error(rw, "Movie not found", http.StatusNotFound)
		return
	}

	reviews, err := rh.service.GetReviewsByMovieID(uint(id))
	if err != nil {
		http.Error(rw, "Server error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(models.MovieRating{
		MovieID:     uint(id),
		AvrRating:   avrRating,
		ReviewCount: len(reviews),
	})
}
