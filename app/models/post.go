package models

import "time"

type Post struct {
    ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID     uint64     `gorm:"index;not null" json:"user_id"`
    Content    string     `gorm:"type:text" json:"content"`
    Anonymous  bool       `gorm:"default:false" json:"anonymous"`
    Visibility string     `gorm:"type:varchar(32);default:'public'" json:"visibility"` // e.g. "public" / "private"
    Images     []Image    `gorm:"foreignKey:PostID" json:"images,omitempty"`
    CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
}