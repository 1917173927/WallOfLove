// Package models 包含所有数据模型定义，负责描述数据库表结构和字段映射关系。
package models

import (
	"time"
)

// Reply 定义回复的数据结构，对应数据库中的回复表。
type Reply struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`      // 回复唯一标识符
	ReviewID  uint64    `json:"review_id"`               // 关联的评论ID
	UserID    uint64    `json:"user_id"`                 // 回复用户的ID
	Content   string    `json:"content"`                 // 回复内容
	CreatedAt time.Time `json:"created_at"`              // 回复创建时间
}