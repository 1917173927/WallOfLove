// Package models 包含所有数据模型定义，负责描述数据库表结构和字段映射关系。
package models

import (
	"time"
)

// Review 定义评论的数据结构，对应数据库中的评论表。
type Review struct {
	ID        uint64      `json:"id" gorm:"primaryKey"`  // 评论唯一标识符
	UserID    uint64      `json:"user_id"`             // 评论用户的ID
	PostID    uint64      `json:"post_id"`             // 关联的帖子ID
	Content   string      `json:"content"`             // 评论内容
	CreatedAt time.Time   `json:"created_at"`          // 评论创建时间
	Replies  []Reply      `json:"replies" gorm:"foreignKey:ReviewID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // 关联的回复列表
}
