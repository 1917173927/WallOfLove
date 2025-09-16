package models

import "time"

type Blacklist struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	BlockedID uint64 `json:"blocked_id"`
	CreatedAt time.Time `json:"created_at"`
}
