package handlers

import (
	"auth-service/internal/service"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (u *UserHandler) Register(rw http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(rw, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if err := u.userService.RegisterUser(data["name"], data["email"], data["password"]); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

func (u *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(rw, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	token, err := u.userService.LoginUser(data["email"], data["password"])
	if err != nil {
		http.Error(rw, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Login successful",
	})
}

func (u *UserHandler) User(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(rw, "Failed to get cookie", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	user, err := u.userService.GetUserFromToken(tokenString)
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(user)

}

func (u *UserHandler) LogOut(rw http.ResponseWriter, r *http.Request) {
	http.SetCookie(rw, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	json.NewEncoder(rw).Encode(map[string]string{
		"message": "Logout Successfully",
	})
}

func (u *UserHandler) ValidateToken(rw http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(rw, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := u.userService.GetUserFromToken(tokenString)
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(user)
}
