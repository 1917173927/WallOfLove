package models

import (
	"time"
)

type Reply struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	ReviewID  uint64    `json:"review_id"`
	UserID    uint64    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}