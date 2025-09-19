package models

import "time"

type User struct {
	ID            uint64    `json:"ID"`
	Username      string    `json:"username"`
	Password      string    `json:"-"`
	Nickname      string    `json:"nickname"`
	AvatarPath    string   `json:"avatar_path"`
	CreatedAt     time.Time `json:"created_at"`
}
