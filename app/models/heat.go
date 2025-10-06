package models

import "time"

type Heat struct {
	ID        uint64    `json:"id"`
	PostID    uint64    `json:"post_id" gorm:"primaryKey"`
	HeatValue uint64    `json:"heat_value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}