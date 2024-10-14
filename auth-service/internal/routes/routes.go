package routes

import (
	"auth-service/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Routes struct {
	r chi.Router
}

func NewRoutes(r chi.Router) *Routes {
	return &Routes{r}
}

func (rt *Routes) SetupRouteUser(userHandler *handlers.UserHandler) {
	rt.r.Post("/register", userHandler.Register)
	rt.r.Post("/login", userHandler.Login)
	rt.r.Get("/user", userHandler.User)
	rt.r.Get("/validate-token", userHandler.ValidateToken)
	rt.r.Post("/logout", userHandler.LogOut)
}
