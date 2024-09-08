package auth

import "github.com/gofrs/uuid"

type User struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"-"`
	Email     string         `json:"email"`
	CreatedAt uuid.Timestamp `json:"createdAt"`
}

type UserOutput struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
