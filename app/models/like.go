package models

import "time"

type Like struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	PostID    uint64    `json:"post_id"`
	ReviewID  uint64    `json:"review_id"`
	CreatedAt time.Time `json:"created_at"`
}