package models

type User struct {
	ID       int    `json:"userID"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	ImageURL string `json:"imageURL"` //头像地址
}

