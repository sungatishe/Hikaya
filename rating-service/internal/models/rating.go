package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	MovieID uint   `json:"movie_id"`
	Rating  int    `json:"rating"`
	Review  string `json:"review"`
}

type MovieRating struct {
	MovieID     uint    `json:"movie_id" gorm:"primaryKey"`
	AvrRating   float64 `json:"avr_rating"`
	ReviewCount int     `json:"review_count"`
}
