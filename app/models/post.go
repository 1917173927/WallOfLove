package models

import "time"

type Post struct {
	ID      int       `json:"ID"`
	UserID int       `json:"UserID"`
	Title string    `json:"Title"`
	Content string    `json:"Content"`
	ImageURL  string    `json:"ImageURL"` //图片地址
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Visible   int       `json:"Visible"` // 公开程度0: public, 1: private,2:DELETED
	Annonymous int       `json:"Anonymous"` // 匿名性0: not anonymous, 1: anonymous
}


