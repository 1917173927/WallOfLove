package models

import "time"

type Post struct {
    ID           uint64     `json:"id"`
    UserID       uint64     `json:"user_id"`
    Content      string     `json:"content"`
    Anonymous    bool       `json:"anonymous"`
    Visibility   bool       `json:"visibility"`
    UserNickname string     `json:"user_nickname"`
    AvatarImageID *uint64   `json:"avatar_image_id"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    Version      uint       `grom:"default:1" json:"version"`
    Images     []Image      `json:"images,omitempty" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
