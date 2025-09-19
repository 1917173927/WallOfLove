package models

import "time"

type Post struct {
    ID           uint64     `json:"id"`
    UserID       uint64     `json:"user_id"`
    Content      string     `json:"content"`
    Anonymous    bool       `json:"anonymous"`
    Visibility   bool       `json:"visibility"`
    UserNickname string     `json:"user_nickname"`
    AvatarPath   string     `json:"avatar_path"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    Images     []Image      `json:"images,omitempty" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
