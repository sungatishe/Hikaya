package repository

import (
	"fmt"
	"gorm.io/gorm"
	"user-list-service/internal/models"
)

type UserListRepository interface {
	CreateUserMovieInList(userList models.UserList) error
	GetUserMovieList(userID uint) ([]models.UserList, error)
	UpdateUserListTypeMovieList(userID, movieID uint, listType string) error
	DeleteFromUserList(userID, movieID uint) error
	isMovieInAnyList(userID, movieID uint) (bool, error)
}

type userListRepository struct {
	db *gorm.DB
}

func NewUserListRepository(db *gorm.DB) UserListRepository {
	return &userListRepository{db}
}

func (u *userListRepository) CreateUserMovieInList(userList models.UserList) error {
	exists, err := u.isMovieInAnyList(userList.UserID, userList.MovieID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("movie is already in one of the lists")
	}

	return u.db.Create(userList).Error
}

func (u *userListRepository) GetUserMovieList(userID uint) ([]models.UserList, error) {
	var userList []models.UserList
	if err := u.db.Where("user_id = ?", userID).Find(&userList).Error; err != nil {
		return nil, err
	}
	return userList, nil
}

func (u *userListRepository) UpdateUserListTypeMovieList(userID, movieID uint, listType string) error {
	exists, err := u.isMovieInAnyList(userID, movieID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("movie not found in any list")
	}

	return u.db.Model(&models.UserList{}).Where("user_id = ? AND movie_id = ?", userID, movieID).
		Update("list_type", listType).Error
}

func (u *userListRepository) DeleteFromUserList(userID, movieID uint) error {
	exists, err := u.isMovieInAnyList(userID, movieID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("movie not found in any list")
	}
	return u.db.Where("user_id = ? AND movie_id = ?", userID, movieID).Delete(&models.UserList{}).Error
}

func (u *userListRepository) isMovieInAnyList(userID, movieID uint) (bool, error) {
	var count int64
	err := u.db.Model(&models.UserList{}).Where("user_id = ? AND movie_id = ?", userID, movieID).Count(&count).Error
	return count > 0, err
}
