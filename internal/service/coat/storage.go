package coat

import (
	"context"
	"errors"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/client/postgresql"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/Henus321/boney-james-go-backend/pkg/utils"
	"github.com/jackc/pgconn"
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

func (s *Storage) GetAll(ctx context.Context) (*[]CoatWithOption, error) {
	const op = "coat.storage.GetAll"

	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetOneByID(ctx context.Context, id string) (*CoatWithOption, error) {
	const op = "coat.storage.GetOneByID"

	query := `SELECT * FROM coat
		LEFT JOIN coat_option 
		ON coat.id = coat_option.coatId
		WHERE coat.id = $1;`

	rows, err := s.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get coats: %w", op, err)
	}
	defer rows.Close()

	var cwo CoatWithOption

	for rows.Next() {
		var (
			id         pgtype.UUID
			colorLabel pgtype.Text
			colorHex   pgtype.Text
			cost       pgtype.Int4
			sizes      pgtype.TextArray
			photoUrls  pgtype.TextArray
			coatId     pgtype.UUID
		)

		err = rows.Scan(
			&cwo.ID,
			&cwo.Name,
			&cwo.Model,
			&cwo.Description,
			&id,
			&colorLabel,
			&colorHex,
			&cost,
			&sizes,
			&photoUrls,
			&coatId,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get coat: %w", op, err)
		}

		if id.Status == pgtype.Present &&
			colorLabel.Status == pgtype.Present &&
			colorHex.Status == pgtype.Present &&
			cost.Status == pgtype.Present &&
			coatId.Status == pgtype.Present &&
			sizes.Status == pgtype.Present &&
			photoUrls.Status == pgtype.Present {
			cwo.CoatOptions = append(cwo.CoatOptions, CoatOption{
				ID:         id.Bytes,
				ColorLabel: colorLabel.String,
				ColorHex:   colorHex.String,
				Cost:       cost.Int,
				Sizes:      utils.FromTextArray(sizes),
				PhotoUrls:  utils.FromTextArray(photoUrls),
				CoatId:     coatId.Bytes,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to get coats: %w", op, err)
	}

	return &cwo, nil
}

// Create по ссылке, можно не возращать т.к это мутирует начальный
func (s *Storage) Create(ctx context.Context, input CreateCoatInput) error {
	q := `
	INSERT INTO coat 
  	  (name)
	VALUES 
   	 ($1)
	RETURNING id
	`

	if err := s.client.QueryRow(ctx, q).Scan(input.Name, input.Description, input.Model); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			s.logger.Error(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
