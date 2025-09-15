package models

import "time"

type User struct {
    ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
    Username      string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
    Nickname      string     `gorm:"type:varchar(100)" json:"nickname,omitempty"`
    AvatarImageID *uint64    `gorm:"index" json:"avatar_image_id,omitempty"` // 可为 NULL
    CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}