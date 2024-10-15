package handlers

import (
	"api-gateway/config"
	"api-gateway/config/cache"
	"api-gateway/internal/rabbitMQ"
	"api-gateway/pgk/httpClient"
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"time"
)

type APIHandlerMovie struct {
	cfg         *config.Config
	rabbitMQ    *rabbitMQ.RabbitMQ
	RedisClient *cache.RedisClient
}

// Movie структура фильма для ответа
type Movie struct {
	Title       string `json:"title"`
	Poster      string `json:"poster"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

func NewAPIHandlerMovie(cfg *config.Config, RedisClient *cache.RedisClient) APIHandlerMovie {
	rmq := rabbitMQ.NewRabbitMQ(cfg.RabbitMQURL)
	return APIHandlerMovie{cfg, rmq, RedisClient}
}

func (a *APIHandlerMovie) GetMovies(rw http.ResponseWriter, r *http.Request) {
	cacheKey := "movies_cache"

	cachedMovies, err := a.RedisClient.GetCache(cacheKey)
	if err != nil {
		http.Error(rw, "Error accessing cache", http.StatusInternalServerError)
		return
	}

	if cachedMovies != "" {
		log.Println("Cache hit, returning cached data")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(cachedMovies))
		return
	}

	log.Println("Cache is nil, fetching from movie service..")

	resp, err := httpClient.GetRequest(a.cfg.MovieServiceURL + "/movies")
	if err != nil {
		http.Error(rw, "Error getting movies", http.StatusInternalServerError)
		return
	}

	err = a.RedisClient.SetCache(cacheKey, string(resp), 10*time.Minute)
	if err != nil {
		log.Printf("Failed to cache movies: %v", err)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (a *APIHandlerMovie) GetMovieByID(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	resp, err := httpClient.GetRequest(a.cfg.MovieServiceURL + "/movies/" + idStr)
	if err != nil {
		http.Error(rw, "Error getting movie", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (a *APIHandlerMovie) UpdateMovie(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error reading request", http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.PutRequest(a.cfg.MovieServiceURL+"/movies/"+idStr, bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		http.Error(rw, "Error updating movie", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (a *APIHandlerMovie) DeleteMovieMQ(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	event := map[string]string{"id": idStr}
	body, err := json.Marshal(event)
	if err != nil {
		http.Error(rw, "Error in marshaling", http.StatusInternalServerError)
		return
	}

	err = a.rabbitMQ.Publish("movie_delete", body)
	if err != nil {
		http.Error(rw, "Error in sending event", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	rw.Write([]byte("event sended"))
}

func (a *APIHandlerMovie) CreateMovieMQ(rw http.ResponseWriter, r *http.Request) {
	var movie Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(rw, "Error in reading request", http.StatusBadRequest)
		return
	}

	if movie.Title == "" || movie.Year == 0 {
		http.Error(rw, "No required fields", http.StatusBadRequest)
		return
	}

	message, err := json.Marshal(movie)
	if err != nil {
		http.Error(rw, "Error in marshalling", http.StatusInternalServerError)
		return
	}

	err = a.rabbitMQ.Publish("movie_create", message)
	if err != nil {
		log.Printf("Error in sending message", err)
		http.Error(rw, "Error in sending event", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	rw.Write([]byte("event create movie sended"))
}
