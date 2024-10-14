package handlers

import (
	"api-gateway/config"
	"api-gateway/pgk/httpClient"
	"fmt"
	"io"
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
	resp, err := httpClient.PostRequestLogin(a.cfg.AuthServiceURL+"/login", r.Body)
	if err != nil {
		log.Println("Error in login: ", err)
		http.Error(rw, fmt.Sprintf("Error login: %v", err), http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}

	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

func (a *APIHandlerAuth) User(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(rw, "Failed to get cookie", http.StatusUnauthorized)
		return
	}

	req, err := http.NewRequest("GET", a.cfg.AuthServiceURL+"/user", nil)
	if err != nil {
		http.Error(rw, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error in getting user data: ", err)
		http.Error(rw, fmt.Sprintf("Error getting user data: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

func (a *APIHandlerAuth) LogOut(rw http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.PostRequestLogin(a.cfg.AuthServiceURL+"/logout", r.Body)
	if err != nil {
		log.Println("Error in logout: ", err)
		http.Error(rw, fmt.Sprintf("Error logout: %v", err), http.StatusInternalServerError)
		return
	}

	// Здесь мы считываем заголовки из ответа auth-service
	for name, values := range resp.Header {
		for _, value := range values {
			rw.Header().Add(name, value)
		}
	}

	// Устанавливаем статус и тело ответа
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, resp.Body)
}
