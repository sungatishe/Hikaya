package service

import (
	"user-list-service/internal/models"
	"user-list-service/internal/repository"
)

type UserListService interface {
	CreateUserMovieInList(userList models.UserList) error
	GetUserMovieList(userID uint) ([]models.UserList, error)
	UpdateUserListTypeMovieList(userID, movieID uint, listType string) error
	DeleteFromUserList(userID, movieID uint) error
	//HandleUserCreated(userData map[string]interface{}) error
}

type userListService struct {
	userRepo repository.UserListRepository
}

func NewUserService(userRepo repository.UserListRepository) UserListService {
	return &userListService{userRepo}
}

func (u *userListService) CreateUserMovieInList(userList models.UserList) error {
	return u.userRepo.CreateUserMovieInList(userList)
}

func (u *userListService) GetUserMovieList(userID uint) ([]models.UserList, error) {
	return u.userRepo.GetUserMovieList(userID)
}

func (u *userListService) UpdateUserListTypeMovieList(userID, movieID uint, listType string) error {
	return u.userRepo.UpdateUserListTypeMovieList(userID, movieID, listType)
}

func (u *userListService) DeleteFromUserList(userID, movieID uint) error {
	return u.userRepo.DeleteFromUserList(userID, movieID)
}

//
//func (u *userListService) HandleUserCreated(userData map[string]interface{}) error {
//	user := &models.User{
//		ID:    int(userData["id"].(float64)), // Преобразование float64 в int
//		Name:  userData["name"].(string),
//		Email: userData["email"].(string),
//	}
//
//	if err := s.repo.CreateUser(user); err != nil {
//		return fmt.Errorf("failed to create user: %w", err)
//	}
//
//	log.Printf("User %s added to user list", user.Name)
//	return nil
//}
