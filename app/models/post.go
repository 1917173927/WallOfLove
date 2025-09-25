// Package models 包含所有数据模型定义，负责描述数据库表结构和字段映射关系。
package models

import "time"

// Post 定义帖子的数据结构，对应数据库中的帖子表。
type Post struct {
    ID           uint64     `json:"id"`                  // 帖子唯一标识符
    IsPublished bool       `json:"is_published"`      // 帖子是否已发布
    UserID       uint64     `json:"user_id"`            // 发帖用户的ID
    Content      string     `json:"content"`            // 帖子内容
    Anonymous    bool       `json:"anonymous"`          // 是否匿名发帖
    Visibility   bool       `json:"visibility"`         // 帖子是否可见
    UserNickname string     `json:"user_nickname"`      // 发帖用户的昵称
    AvatarPath   string     `json:"avatar_path"`        // 发帖用户的头像路径
    ScheduledAt  *time.Time `json:"scheduled_at"`       // 帖子计划发布时间，可为空
    CreatedAt    time.Time  `json:"created_at"`         // 帖子创建时间
    UpdatedAt    time.Time  `json:"updated_at"`         // 帖子最后更新时间
    Images     []Image      `json:"images,omitempty" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // 关联的图片列表
    LikeCount    int64      `json:"like_count"  gorm:"-"`
	LikedByMe    bool       `json:"liked_by_me" gorm:"-"`
}
