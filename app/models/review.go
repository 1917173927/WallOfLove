package models

import (
	"time"
)

type Review struct {
	ID        uint64      `json:"id" gorm:"primaryKey"`
	UserID    uint64      `json:"user_id"`
	PostID    uint64      `json:"post_id"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	Replies  []Reply   `json:"replies" gorm:"foreignKey:ReviewID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
