package models

import "time"

type Document struct {
	Id      int       `json:"-" db:"id"`
	Name    string    `json:"name" binding:"required"`
	Mime    string    `json:"mime" binding:"required"`
	Token   string    `json:"token" binding:"required"`
	File    bool      `json:"file" binding:"required"`
	Public  bool      `json:"public" binding:"required"`
	Grant   []string  `json:"grant" binding:"required"`
	Created time.Time `json:"created"`
}
