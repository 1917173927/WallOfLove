
package models

import (
	"time"
)


type Reply struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`      // 回复唯一标识符
	ReviewID  uint64    `json:"review_id"`               // 关联的评论ID
	UserID    uint64    `json:"user_id"`                 // 回复用户的ID
	Content   string    `json:"content"`                 // 回复内容
	CreatedAt time.Time `json:"created_at"`              // 回复创建时间
}

type ReplyWithNickname struct {
	Reply
	Nickname   string `json:"nickname"`
	AvatarPath string `json:"avatarPath"`
}