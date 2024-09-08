package auth

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/client/postgresql"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/jackc/pgtype"
)

type Storage struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewStorage(client postgresql.Client, logger *logging.Logger) *Storage {
	return &Storage{
		client: client,
		logger: logger,
	}
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*UserOutput, error) {
	const op = "auth.storage.GetUserByID"

	query := `SELECT id, username, email FROM users WHERE id = $1`

	row := s.client.QueryRow(ctx, query, id)

	var user UserOutput

	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const op = "auth.storage.GetUserByEmail"

	query := `SELECT id, username, email, password, createdAt FROM users WHERE email = $1;`

	row := s.client.QueryRow(ctx, query, email)

	var user User

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) CreateUser(ctx context.Context, input *UserRegisterInput) (*UserOutput, error) {
	const op = "auth.storage.CreateUser"

	query := `
		INSERT INTO users 
		  (username, email, password)
		VALUES 
		 ($1, $2, $3)
		RETURNING id
	`

	var (
		newId pgtype.UUID
	)

	err := s.client.QueryRow(ctx, query, input.Username, input.Email, input.Password).Scan(&newId)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create user: %w", op, err)
	}

	user := UserOutput{
		ID:       newId.Bytes,
		Username: input.Username,
		Email:    input.Email,
	}

	return &user, nil
}
