package handlers

import (
	"encoding/json"
	"net/http"
	"user-list-service/internal/models"
	"user-list-service/internal/service"
)

type UserListHandler struct {
	movieService service.UserService
}

func NewUserListHandler(movieService service.UserService) *UserListHandler {
	return &UserListHandler{movieService: movieService}
}

func (u *UserListHandler) AddMovieToUserList(rw http.ResponseWriter, r *http.Request) {
	var userMovie models.UserList
	if err := json.NewDecoder(r.Body).Decode(&userMovie); err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := u.movieService.CreateUserMovieInList(userMovie); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(userMovie)
}
