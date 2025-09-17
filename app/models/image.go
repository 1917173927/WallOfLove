package models

import "time"

type Image struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	PostID     *uint64    `json:"post_id"`
	IsAvatar    bool      `json:"is_avatar"`
	FilePath    string    `json:"file_path"`
	Size        int64     `json:"size"`
	Checksum    string    `json:"checksum"`
	CreatedAt   time.Time `json:"created_at"`
}