package models

type View struct {
	ID        uint64   `json:"id"`
	UserID    uint64   `json:"user_id"`
	PostID    uint64   `json:"post_id"`
	CreatedAt string `json:"created_at"`
}