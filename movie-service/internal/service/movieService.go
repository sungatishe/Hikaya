package service

import (
	"fmt"
	"movie-service/internal/models"
	"movie-service/internal/repository"
	"strconv"
)

type MovieService interface {
	CreateMovie(movie *models.Movie) error
	GetMovieById(id uint) (*models.Movie, error)
	GetAllMovies() ([]models.Movie, error)
	UpdateMovie(movie *models.Movie) error
	DeleteMovie(id uint) error
	IndexMovie(movieID string, movieData models.Movie) error
	SearchMovies(query string) ([]map[string]interface{}, error)
	IndexAllMovies() error
}

type movieService struct {
	movieRepo  repository.MovieRepository
	searchRepo repository.ElasticSearchRepository
}

func NewMovieService(movieRepo repository.MovieRepository, searchRepo repository.ElasticSearchRepository) MovieService {
	return movieService{movieRepo, searchRepo}
}

func (m movieService) CreateMovie(movie *models.Movie) error {
	return m.movieRepo.CreateMovie(movie)
}

func (m movieService) GetMovieById(id uint) (*models.Movie, error) {
	return m.movieRepo.GetMovieById(id)
}

func (m movieService) GetAllMovies() ([]models.Movie, error) {
	return m.movieRepo.GetAllMovies()
}

func (m movieService) UpdateMovie(movie *models.Movie) error {
	return m.movieRepo.UpdateMovie(movie)
}

func (m movieService) DeleteMovie(id uint) error {
	return m.movieRepo.DeleteMovie(id)
}

func (m movieService) IndexMovie(movieID string, movieData models.Movie) error {
	return m.searchRepo.IndexMovie(movieID, movieData)
}

func (m movieService) SearchMovies(query string) ([]map[string]interface{}, error) {
	return m.searchRepo.SearchMovies(query)
}

func (m movieService) IndexAllMovies() error {
	movies, err := m.movieRepo.GetAllMovies()
	if err != nil {
		return fmt.Errorf("could not retrieve movies from database: %w", err)
	}
	for _, movie := range movies {
		err := m.searchRepo.IndexMovie(strconv.Itoa(int(movie.ID)), movie)
		if err != nil {
			return fmt.Errorf("could not index movie %d: %w", movie.ID, err)
		}
	}
	return nil
}
