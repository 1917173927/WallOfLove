package models

import "time"

type User struct {
    ID             uint64     `json:"id"`
    Username       string     `json:"username"`
    Password       string     `json:"-"`
    Nickname       string     `json:"nickname"`
    AvatarImageID *uint64     `json:"avatar_image_id"`
    CreatedAt      time.Time  `json:"created_at"`
    Version        uint       `grom:"default:1" json:"version"`
}