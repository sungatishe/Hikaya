package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"user-list-service/internal/models"
	"user-list-service/internal/service"
)

type UserListHandler struct {
	movieService service.UserListService
}

func NewUserListHandler(movieService service.UserListService) *UserListHandler {
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

func (u *UserListHandler) GetUsersMovieList(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	list, err := u.movieService.GetUserMovieList(uint(id))
	if err != nil {
		http.Error(rw, "List not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(rw).Encode(list)
}

func (u *UserListHandler) UpdateMovieListType(rw http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	movieIDStr := chi.URLParam(r, "movieID")
	newListType := r.URL.Query().Get("listType")

	userID, err := strconv.ParseUint(userIDStr, 10, 0)
	if err != nil {
		http.Error(rw, "Invalid user ID", http.StatusBadRequest)
		return
	}

	movieID, err := strconv.ParseUint(movieIDStr, 10, 0)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	err = u.movieService.UpdateUserListTypeMovieList(uint(userID), uint(movieID), newListType)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Movie list type updated successfully"))
}

func (u *UserListHandler) DeleteMovieFromUserList(rw http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	movieIDStr := chi.URLParam(r, "movieID")

	userID, err := strconv.ParseUint(userIDStr, 10, 0)
	if err != nil {
		http.Error(rw, "Invalid user ID", http.StatusBadRequest)
		return
	}

	movieID, err := strconv.ParseUint(movieIDStr, 10, 0)
	if err != nil {
		http.Error(rw, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	err = u.movieService.DeleteFromUserList(uint(userID), uint(movieID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Movie deleted from user's list successfully"))
}
