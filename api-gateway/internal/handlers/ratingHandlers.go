package handlers

import (
	"api-gateway/config"
	"api-gateway/pgk/httpClient"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type APIHandlerRating struct {
	cfg *config.Config
}

func NewAPIHandlerRating(cfg *config.Config) APIHandlerRating {
	return APIHandlerRating{cfg: cfg}
}

func (a *APIHandlerRating) CreateMovieReview(rw http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.PostRequest(a.cfg.RatingServiceURL+"/reviews", r.Body)
	if err != nil {
		http.Error(rw, "Error reading request", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write(resp)
}

func (a *APIHandlerRating) GetMovieReviews(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	url := fmt.Sprintf("/movie/%s/reviews", idStr)
	resp, err := httpClient.GetRequest(a.cfg.RatingServiceURL + url)
	if err != nil {
		http.Error(rw, "Error getting movie reviews", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (a *APIHandlerRating) GetMovieRating(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	url := fmt.Sprintf("/movie/%s/rating", idStr)
	resp, err := httpClient.GetRequest(a.cfg.RatingServiceURL + url)
	if err != nil {
		http.Error(rw, "Error getting movie rating", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}
