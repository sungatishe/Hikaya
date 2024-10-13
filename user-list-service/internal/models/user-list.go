package models

import "gorm.io/gorm"

type ListType string

const (
	WatchedList   ListType = "watched"
	PlannedList   ListType = "planned"
	AbandonedList ListType = "abandoned"
)

type UserList struct {
	gorm.Model
	UserID   uint     `json:"user_id"`
	MovieID  uint     `json:"movie_id"`
	ListType ListType `json:"list_type"`
}
