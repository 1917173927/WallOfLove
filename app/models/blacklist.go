package models

import "time"

type Blacklist struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	BlackUserID uint64    `json:"black_user_id"`
	CreatedAt   time.Time `json:"created_at"`
}