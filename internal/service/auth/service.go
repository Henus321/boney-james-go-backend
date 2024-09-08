package auth

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	storage *Storage
}

type Service interface {
	GetUserByID(ctx context.Context, id string) (*UserOutput, error)
	LoginUser(ctx context.Context, input *UserLoginInput) (token string, err error)
	RegisterUser(ctx context.Context, input *UserRegisterInput) error
	CreateJWT(secret []byte, id string) (string, error)
	HashPassword(password string) (string, error)
	ComparePasswords(hashed string, plain []byte) bool
}

func NewService(storage *Storage) Service {
	return &service{storage: storage}
}

func (s service) GetUserByID(ctx context.Context, id string) (*UserOutput, error) {
	return s.storage.GetUserByID(ctx, id)
}

func (s service) LoginUser(ctx context.Context, input *UserLoginInput) (token string, err error) {
	const op = "auth.service.LoginUser"

	user, err := s.storage.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return "", fmt.Errorf("%s: invalid credentials", op)
	}

	if !s.ComparePasswords(user.Password, []byte(input.Password)) {
		return "", fmt.Errorf("%s: invalid credentials", op)
	}

	secret := []byte(config.GetConfig().JWT.Secret)
	token, err = s.CreateJWT(secret, user.ID)
	if err != nil {
		return "", fmt.Errorf("%s: cant create token", op)
	}

	return token, nil
}

func (s service) RegisterUser(ctx context.Context, input *UserRegisterInput) error {
	const op = "auth.service.RegisterUser"

	user, _ := s.storage.GetUserByEmail(ctx, input.Email)
	if user != nil {
		return fmt.Errorf("%s: email already in use", op)
	}

	hashedPassword, err := s.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("%s: failed to hash password: %w", op, err)
	}
	input.Password = hashedPassword

	err = s.storage.CreateUser(ctx, input)
	if err != nil {
		return fmt.Errorf("%s: failed to create user: %w", op, err)
	}

	return nil
}

func (s service) CreateJWT(secret []byte, id string) (string, error) {
	expiration := time.Second * time.Duration(3600*24*7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    id,
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (s service) HashPassword(password string) (string, error) {
	const op = "auth.storage.HashPassword"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	return string(hash), nil
}

func (s service) ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)

	return err == nil
}
