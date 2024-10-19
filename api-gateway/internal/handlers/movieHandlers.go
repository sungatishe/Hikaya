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

// GetMovies godoc
// @Summary Get all movies
// @Description Получить список всех фильмов
// @Tags movies
// @Produce json
// @Success 200 {array} Movie
// @Failure 500 {string} string "Error getting movies"
// @Router /movies [get]
// GetMovies retrieves all movies from the cache or fetches them from the movie service if not cached.
// It returns a list of movies in JSON format.
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

	var movies interface{}
	if err := json.Unmarshal(resp, &movies); err != nil {
		http.Error(rw, "Error parsing movie data", http.StatusInternalServerError)
		return
	}

	formattedMovies, err := json.MarshalIndent(movies, "", "    ") // 4 пробела для отступа
	if err != nil {
		http.Error(rw, "Error formatting movie data", http.StatusInternalServerError)
		return
	}

	err = a.RedisClient.SetCache(cacheKey, string(resp), 10*time.Minute)
	if err != nil {
		log.Printf("Failed to cache movies: %v", err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(formattedMovies)
}

// GetMovieByID godoc
// @Summary Get movie by ID
// @Description Получить информацию о фильме по его ID
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} Movie
// @Failure 404 {string} string "Movie not found"
// @Router /movies/{id} [get]
// GetMovieByID retrieves information about a specific movie identified by its ID.
// It responds with the movie details in JSON format.
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

// UpdateMovie godoc
// @Summary Update a movie
// @Description Обновить информацию о фильме по его ID
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Param movie body Movie true "Movie details"
// @Success 200 {string} string "Movie updated successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Error updating movie"
// @Router /movies/{id} [put]
// UpdateMovie updates the details of an existing movie based on the provided ID and request body.
// It expects a JSON object with the updated movie details.
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

// DeleteMovieMQ godoc
// @Summary Delete a movie
// @Description Удалить фильм по его ID, отправив событие в очередь RabbitMQ
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Success 202 {string} string "Event sent successfully"
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Error in sending event"
// @Router /movies/{id} [delete]
// DeleteMovieMQ sends a delete event for a movie identified by its ID to RabbitMQ.
// This triggers the deletion process in the appropriate microservice.
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

// CreateMovieMQ godoc
// @Summary Create a new movie
// @Description Создать новый фильм, отправив событие в очередь RabbitMQ
// @Tags movies
// @Produce json
// @Param movie body Movie true "Movie details"
// @Success 202 {string} string "Event sent successfully"
// @Failure 400 {string} string "No required fields"
// @Failure 500 {string} string "Error in sending event"
// @Router /movies [post]
// CreateMovieMQ creates a new movie by sending a create event to RabbitMQ.
// It expects a JSON object with the movie details in the request body.
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

func (a *APIHandlerMovie) SearchMovies(rw http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(rw, "Invalid query", http.StatusBadRequest)
		return
	}

	resp, err := httpClient.GetRequest(a.cfg.MovieServiceURL + "/movies/search?q=" + query)
	if err != nil {
		http.Error(rw, "Error getting movies from search", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}
