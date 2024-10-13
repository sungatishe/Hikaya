package service

import (
	"auth-service/internal/models"
	"auth-service/internal/repository"
	rabbitmq "auth-service/internal/rmq"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService interface {
	RegisterUser(name, email, password string) error
	LoginUser(email, password string) (string, error)
	GetUserFromToken(tokenString string) (*models.User, error)
}

type userService struct {
	userRepo  repository.UserRepository
	secretKey string
	rabbit    *rabbitmq.RabbitMQ
}

func NewUserService(userRepo repository.UserRepository, secretKey string, rabbit *rabbitmq.RabbitMQ) UserService {
	return userService{userRepo: userRepo, secretKey: secretKey, rabbit: rabbit}
}

func (u userService) RegisterUser(name, email, password string) error {
	if _, err := u.userRepo.GetUserByEmail(email); err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash")
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.CreateUser(&user); err != nil {
		return errors.New("failed to create User")
	}

	// Публикация события в RabbitMQ
	event := map[string]interface{}{
		"eventType": "UserCreated",
		"data":      user,
	}
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	if err := u.rabbit.PublishMessage(string(eventData)); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	return nil
}

func (u userService) LoginUser(email, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Invalid credentials")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(u.secretKey))
	if err != nil {
		fmt.Println("Error generating token:", err)
		return "", errors.New("Error token")
	}

	return token, nil
}

func (u userService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.secretKey), nil
	})

	if err != nil {
		return nil, errors.New("Unauthorized")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, errors.New("Failed to parse claims")
	}

	id, _ := strconv.Atoi((*claims)["sub"].(string))
	user, err := u.userRepo.GetUserById(uint(id))
	if err != nil {
		return nil, errors.New("Failed to get ID")
	}
	return user, nil
}
