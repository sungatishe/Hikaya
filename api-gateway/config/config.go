package config

import "os"

type Config struct {
	Port            string
	AuthServiceURL  string
	MovieServiceURL string
	JWTSecret       string
	RabbitMQURL     string
}

func LoadConfig() *Config {
	return &Config{
		Port:            os.Getenv("PORT"),
		AuthServiceURL:  os.Getenv("AUTH_SERVICE_URL"),
		MovieServiceURL: os.Getenv("MOVIE_SERVICE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		RabbitMQURL:     os.Getenv("RABBITMQURL"),
	}
}
