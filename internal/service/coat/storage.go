package coat

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/client/postgresql"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/Henus321/boney-james-go-backend/pkg/utils"
	"github.com/gofrs/uuid"
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

	query := `SELECT 
        	coat.id,
			name,
			model,
			description,
			coat_option.id as optionId,
			colorLabel,
			colorHex,
			cost,
			sizes,
			photoUrls,
			coatId
       	FROM coat
		LEFT JOIN coat_option 
		ON coat.id = coat_option.coatId`

	rows, err := s.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get coats: %w", op, err)
	}
	defer rows.Close()

	var cwos []CoatWithOption

	for rows.Next() {
		var (
			id          pgtype.UUID
			name        pgtype.Text
			model       pgtype.Text
			description pgtype.Text
			optionId    pgtype.UUID
			colorLabel  pgtype.Text
			colorHex    pgtype.Text
			cost        pgtype.Int4
			sizes       pgtype.TextArray
			photoUrls   pgtype.TextArray
			coatId      pgtype.UUID
		)

		err = rows.Scan(
			&id,
			&name,
			&model,
			&description,
			&optionId,
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

		index := IndexOfCoat(cwos, id.Bytes)

		if id.Status == pgtype.Present &&
			name.Status == pgtype.Present &&
			model.Status == pgtype.Present &&
			description.Status == pgtype.Present &&
			optionId.Status == pgtype.Present &&
			colorLabel.Status == pgtype.Present &&
			colorHex.Status == pgtype.Present &&
			cost.Status == pgtype.Present &&
			coatId.Status == pgtype.Present &&
			sizes.Status == pgtype.Present &&
			photoUrls.Status == pgtype.Present {
			if index == -1 {
				cwos = append(cwos, CoatWithOption{
					ID:          id.Bytes,
					Name:        name.String,
					Model:       model.String,
					Description: description.String,
					CoatOptions: []CoatOption{{
						ID:         optionId.Bytes,
						ColorLabel: colorLabel.String,
						ColorHex:   colorHex.String,
						Cost:       cost.Int,
						Sizes:      utils.FromTextArray(sizes),
						PhotoUrls:  utils.FromTextArray(photoUrls),
						CoatId:     coatId.Bytes,
					}},
				})
			} else {
				cwos[index].CoatOptions = append(cwos[index].CoatOptions, CoatOption{
					ID:         optionId.Bytes,
					ColorLabel: colorLabel.String,
					ColorHex:   colorHex.String,
					Cost:       cost.Int,
					Sizes:      utils.FromTextArray(sizes),
					PhotoUrls:  utils.FromTextArray(photoUrls),
					CoatId:     coatId.Bytes,
				})
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to get coats: %w", op, err)
	}

	return &cwos, nil
}

func (s *Storage) GetOneByID(ctx context.Context, id string) (*CoatWithOption, error) {
	const op = "coat.storage.GetOneByID"

	query := `SELECT 
    		coat.id,
			name,
			model,
			description,
			coat_option.id as optionId,
			colorLabel,
			colorHex,
			cost,
			sizes,
			photoUrls,
			coatId
    	FROM coat
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
			optionId   pgtype.UUID
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
			&optionId,
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

		if optionId.Status == pgtype.Present &&
			colorLabel.Status == pgtype.Present &&
			colorHex.Status == pgtype.Present &&
			cost.Status == pgtype.Present &&
			coatId.Status == pgtype.Present &&
			sizes.Status == pgtype.Present &&
			photoUrls.Status == pgtype.Present {
			cwo.CoatOptions = append(cwo.CoatOptions, CoatOption{
				ID:         optionId.Bytes,
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
	const op = "coat.storage.Create"

	query := `
		INSERT INTO coat 
		  (model, name, description)
		VALUES 
		 ($1, $2, $3)
		RETURNING id
	`

	var (
		newId pgtype.UUID
	)

	err := s.client.QueryRow(ctx, query, input.Model, input.Name, input.Description).Scan(&newId)
	if err != nil {
		return fmt.Errorf("%s: failed to get coats: %w", op, err)
	}

	// ??? returning как правильно вернуть newId? RETURNING id
	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	const op = "coat.storage.Delete"

	query := `
		DELETE FROM coat 
		WHERE id = $1
		RETURNING id
	`

	var (
		oldId pgtype.UUID
	)

	err := s.client.QueryRow(ctx, query, id).Scan(&oldId)
	if err != nil {
		return fmt.Errorf("%s: failed to get coats: %w", op, err)
	}

	// ??? returning как правильно вернуть oldId? RETURNING id
	return nil
}

// IndexOfCoat ??? Куда такие хелперы с конкретным значением
func IndexOfCoat(coats []CoatWithOption, coatId uuid.UUID) int {
	for i, v := range coats {
		if v.ID == coatId {
			return i
		}
	}
	return -1
}
