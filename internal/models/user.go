package models

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"-"`
}
