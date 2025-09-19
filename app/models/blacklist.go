// Package models 包含所有数据模型定义，负责描述数据库表结构和字段映射关系。
package models

import "time"

// Blacklist 定义黑名单的数据结构，对应数据库中的黑名单表。
type Blacklist struct {
	ID        uint64    `json:"id"`                  // 黑名单记录唯一标识符
	UserID    uint64    `json:"user_id"`            // 操作用户的ID
	BlockedID uint64    `json:"blocked_id"`         // 被拉黑用户的ID
	CreatedAt time.Time `json:"created_at"`        // 记录创建时间
}
