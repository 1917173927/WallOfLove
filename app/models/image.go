package models

import "time"

type Image struct {
	ID          uint64    `json:"id"`
	OwnerID     uint64    `json:"owner_id"`
	PostID     *uint64    `json:"post_id"`
	IsAvatar    bool      `json:"is_avatar"`
	FilePath    string    `json:"file_path"`
	ThumbPath   string    `json:"thumb_path"`//缩略图路径
	Mime        string    `json:"mime"`//图片类型
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Size        int64     `json:"size"`
	Checksum    string    `json:"checksum"`
	OrderIndex  int       `json:"order_index"`//图片顺序
	CreatedAt   time.Time `json:"created_at"`
}