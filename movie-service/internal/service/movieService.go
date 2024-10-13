package service

import (
	"movie-service/internal/models"
	"movie-service/internal/repository"
)

type MovieService interface {
	CreateMovie(movie *models.Movie) error
	GetMovieById(id uint) (*models.Movie, error)
	GetAllMovies() ([]models.Movie, error)
	UpdateMovie(movie *models.Movie) error
	DeleteMovie(id uint) error
}

type movieService struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return movieService{movieRepo}
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
