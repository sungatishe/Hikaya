package repository

import (
	"gorm.io/gorm"
	"movie-service/internal/models"
)

type MovieRepository interface {
	CreateMovie(movie *models.Movie) error
	GetMovieById(id uint) (*models.Movie, error)
	GetAllMovies() ([]models.Movie, error)
	UpdateMovie(movie *models.Movie) error
	DeleteMovie(id uint) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return movieRepository{db}
}

func (m movieRepository) CreateMovie(movie *models.Movie) error {
	return m.db.Create(movie).Error
}

func (m movieRepository) GetMovieById(id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := m.db.First(&movie, id).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (m movieRepository) GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie
	if err := m.db.Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (m movieRepository) UpdateMovie(movie *models.Movie) error {
	return m.db.Save(movie).Error
}

func (m movieRepository) DeleteMovie(id uint) error {
	return m.db.Delete(&models.Movie{}, id).Error
}
