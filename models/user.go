package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
}
