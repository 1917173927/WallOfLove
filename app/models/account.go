// Package models 包含所有数据模型定义，负责描述数据库表结构和字段映射关系。
package models

import "time"

// User 定义用户账户的数据结构，对应数据库中的用户表。
type User struct {
	ID            uint64    `json:"ID"`                  // 用户唯一标识符
	Username      string    `json:"username"`            // 用户名，用于登录和显示
	Password      string    `json:"-"`                   // 密码，JSON序列化时忽略
	Nickname      string    `json:"nickname"`            // 用户昵称，用于显示
	AvatarPath    string   `json:"avatar_path"`          // 用户头像路径
	CreatedAt     time.Time `json:"created_at"`          // 用户创建时间
}
