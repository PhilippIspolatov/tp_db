package models

type User struct {
	Nickname string `json:"nickname"`
	Email string `json:"email"`
	FullName string `json:"fullname"`
	About string `json:"about"`
}
