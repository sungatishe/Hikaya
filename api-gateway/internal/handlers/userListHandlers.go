package handlers

import (
	"api-gateway/config"
	"api-gateway/pgk/httpClient"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type APIHandlerUserList struct {
	cfg *config.Config
}

func NewAPIHandlerUserList(cfg *config.Config) APIHandlerUserList {
	return APIHandlerUserList{cfg}
}

func (a *APIHandlerUserList) GetUsersList(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	resp, err := httpClient.GetRequest(a.cfg.UserListServiceURL + "/userList/" + idStr)
	if err != nil {
		http.Error(rw, "Error getting user's list", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (a *APIHandlerUserList) AddMovieToUsersList(rw http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.PostRequest(a.cfg.UserListServiceURL+"/userList", r.Body)
	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write(resp)
}

func (a *APIHandlerUserList) UpdateMovieListType(rw http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	movieIDStr := chi.URLParam(r, "movieID")
	newListType := r.URL.Query().Get("listType")
	// for example: userList/3/6?listType=abandoned
	url := fmt.Sprintf("%s/%s?listType=%s", userIDStr, movieIDStr, newListType)

	resp, err := httpClient.PutRequest(a.cfg.UserListServiceURL+"/userList/"+url, r.Body)
	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (a *APIHandlerUserList) DeleteMovieFromUserList(rw http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	movieIDStr := chi.URLParam(r, "movieID")
	url := fmt.Sprintf("%s/%s", userIDStr, movieIDStr)

	resp, err := httpClient.DeleteRequest(a.cfg.UserListServiceURL + "/userList/" + url)
	if err != nil {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}
