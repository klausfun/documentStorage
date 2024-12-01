package models

type User struct {
	Id       int    `json:"-"`
	Token    string `json:"token" binding:"required"`
	Login    string `json:"login" binding:"required" db:"login"`
	Password string `json:"password" binding:"required"`
}
