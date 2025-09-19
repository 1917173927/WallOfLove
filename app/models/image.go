// Package models 包含所有数据模型定义，负责描述数据库表结构和字段映射关系。
package models

import "time"

// Image 定义图片的数据结构，对应数据库中的图片表。
type Image struct {
	ID          uint64    `json:"id"`                  // 图片唯一标识符
	UserID      uint64    `json:"user_id"`            // 上传图片的用户ID
	PostID      *uint64   `json:"post_id"`            // 关联的帖子ID，可为空
	IsAvatar    bool      `json:"is_avatar"`          // 标记是否为头像图片
	FilePath    string    `json:"file_path"`          // 图片文件存储路径
	Size        int64     `json:"size"`               // 图片文件大小（字节）
	Checksum    string    `json:"checksum"`          // 图片文件的校验和
	CreatedAt   time.Time `json:"created_at"`        // 图片上传时间
}
