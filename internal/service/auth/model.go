package auth

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserOutput struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
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
