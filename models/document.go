package models

type Document struct {
	Id     int      `json:"-"`
	Name   string   `json:"name" binding:"required"`
	Mime   string   `json:"mime" binding:"required"`
	Token  string   `json:"token" binding:"required"`
	File   bool     `json:"file" binding:"required"`
	Public bool     `json:"public" binding:"required"`
	Grant  []string `json:"grant" binding:"required"`
}

type GetDocsInput struct {
	Id    int     `json:"-"`
	Token string  `json:"token" binding:"required"`
	Login *string `json:"login"`
	Key   string  `json:"key" binding:"required"`
	Value string  `json:"value" binding:"required"`
	Limit int     `json:"limit" binding:"required"`
}

type GetDocsResp struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Mime    string   `json:"mime"`
	File    bool     `json:"file"`
	Public  bool     `json:"public"`
	Created string   `json:"created"`
	Grant   []string `json:"grant"`
}

type GetDoc struct {
	File     []byte
	MimeType string
	JSON     string
	IsFile   bool
}
