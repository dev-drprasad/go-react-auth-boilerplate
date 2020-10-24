package model

type User struct {
	ID uint `json:"id"`

	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
