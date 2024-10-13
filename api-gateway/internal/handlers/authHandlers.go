package handlers

import (
	"api-gateway/config"
	"api-gateway/pgk/httpClient"
	"fmt"
	"log"
	"net/http"
)

type APIHandlerAuth struct {
	cfg *config.Config
}

func NewAPIHandlerAuth(cfg *config.Config) APIHandlerAuth {
	return APIHandlerAuth{cfg}
}

func (a *APIHandlerAuth) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.PostRequest(a.cfg.AuthServiceURL+"/register", r.Body)
	if err != nil {
		log.Println("Error in register: ", err)
		http.Error(rw, "Error in register", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (a *APIHandlerAuth) LoginUser(rw http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.PostRequest(a.cfg.AuthServiceURL+"/login", r.Body)
	if err != nil {
		log.Println("Error in login: ", err)
		http.Error(rw, fmt.Sprintf("Error login: ", err), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}
