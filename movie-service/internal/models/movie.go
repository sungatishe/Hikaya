package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Poster      string `json:"poster"`
	Year        int    `json:"year"`
}
