package models

import "time"

type Post struct {
    ID           uint64     `json:"id"`
    UserID       uint64     `json:"user_id"`
    Content      string     `json:"content"`
    Anonymous    bool       `json:"anonymous"`
    Visibility   string     `json:"visibility"` // e.g. "public" / "private"
    Images     []Image      `json:"images,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdateAt     time.Time  `json:"update_at"`
}