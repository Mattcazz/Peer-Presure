package types

import (
	"time"
)

type UserStore interface {
	CreateUser(User) error
	GetUserByEmail(string) (*User, error)
	GetUserById(int) (*User, error)
}

type RegisterUserPayload struct {
	UserName string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"-"`

	CreatedAt time.Time `json:"created_at"`
}
