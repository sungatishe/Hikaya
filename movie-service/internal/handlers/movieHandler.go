package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"movie-service/internal/models"
	"movie-service/internal/service"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	movieService service.MovieService
}

func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

func (m *MovieHandler) CreateMovie(rw http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := m.movieService.CreateMovie(&movie); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(movie)
}

func (m *MovieHandler) GetMovieById(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := m.movieService.GetMovieById(uint(id))
	if err != nil {
		http.Error(rw, "Movie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(rw).Encode(movie)
}

func (m *MovieHandler) GetAllMovies(rw http.ResponseWriter, r *http.Request) {
	movies, err := m.movieService.GetAllMovies()
	if err != nil {
		http.Error(rw, "Failed to retrieve movies", http.StatusBadRequest)
		return
	}

	json.NewEncoder(rw).Encode(movies)

}

func (m *MovieHandler) UpdateMovie(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(rw, "Invalid input", http.StatusBadRequest)
		return
	}
	movie.ID = uint(id)

	if err := m.movieService.UpdateMovie(&movie); err != nil {
		http.Error(rw, "Failed to update movie", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(movie)
}

func (m *MovieHandler) DeleteMovie(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	if err := m.movieService.DeleteMovie(uint(id)); err != nil {
		http.Error(rw, "Failed to delete movie", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (m *MovieHandler) HandleDeleteMovieEvent(body []byte) {
	var event map[string]string
	err := json.Unmarshal(body, &event)
	if err != nil {
		log.Printf("Ошибка при десериализации события: %s", err)
		return
	}

	movieID := event["id"]
	id, err := strconv.Atoi(movieID)
	if err != nil {
		log.Printf("Ошибка при десериализации события: %s", err)
		return
	}

	err = m.movieService.DeleteMovie(uint(id)) // вызов метода удаления фильма
	if err != nil {
		log.Printf("Ошибка при удалении фильма с ID %s: %s", movieID, err)
		return
	}

	log.Printf("Фильм с ID %s был успешно удален", movieID)
}

func (m *MovieHandler) HandleCreateMovieEvent(body []byte) {
	var movie models.Movie
	err := json.Unmarshal(body, &movie)
	if err != nil {
		log.Printf("Ошибка при десериализации события: %s", err)
		return
	}

	err = m.movieService.CreateMovie(&movie) // вызов метода удаления фильма
	if err != nil {
		log.Printf("Error creating movie %s", err)
		return
	}

	log.Printf("Created film")
}
