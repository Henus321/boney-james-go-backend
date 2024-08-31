package coat

import (
	"context"
	"errors"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/client/postgresql"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"github.com/jackc/pgconn"
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

func (s *Storage) GetAll(ctx context.Context) ([]CoatWithOption, error) {
	const op = "coat.storage.GetAll"

	q := `SELECT * FROM coat
		LEFT JOIN coat_option 
		ON coat.id = coat_option.coatid
	`

	rows, err := s.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	coats := make([]CoatWithOption, 0)

	for rows.Next() {
		var cwo CoatWithOption

		err = rows.Scan(&cwo.ID, &cwo.Name, &cwo.Model, &cwo.Description)
		if err != nil {
			return nil, err
		}

		coats = append(coats, cwo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return coats, nil
}

func (s *Storage) GetOneByID(ctx context.Context, id string) (Coat, error) {
	//TODO implement me
	panic("implement me")
}

// Create по ссылке, можно не возращать т.к это мутирует начальный
func (s *Storage) Create(ctx context.Context, dto CreateCoatDTO) error {
	q := `
	INSERT INTO coat 
  	  (name)
	VALUES 
   	 ($1)
	RETURNING id
	`

	if err := s.client.QueryRow(ctx, q).Scan(dto.Name, dto.Description, dto.Model); err != nil {
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

	//TODO implement me
	panic("implement me")
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

//func (r *repository) GetOneByID(ctx context.Context, id string) (coat *Coat, err error) {
//	const op = "coat.storage.GetOneByID"
//
//	oid, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		d.logger.Error(op, err)
//		return coat, fmt.Errorf("%s: failed to convert hex to objectid: %w", op, err)
//	}
//	filter := bson.M{"_id": oid}
//	result := d.collection.FindOne(ctx, filter)
//	if result.Err() != nil {
//		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
//			d.logger.Error(op, err)
//			return coat, fmt.Errorf("%s: failed to find coat by id: %w", op, err)
//		}
//
//		d.logger.Error(op, err)
//		return coat, fmt.Errorf("%s: failed to find coat by id: %w", op, err)
//	}
//
//	if err = result.Decode(&coat); err != nil {
//		d.logger.Error(op, err)
//		return coat, fmt.Errorf("%s: failed to decode coat from DB: %w", op, err)
//	}
//
//	d.logger.Infof("%s: success", op)
//	return coat, nil
//}
//
//func (r *repository) Create(ctx context.Context, dto *CreateCoatDTO) error {
//	const op = "coat.storage.Create"
//
//	collection := d.collection.Database().Collection("coat")
//
//	_, err := collection.InsertOne(ctx, dto)
//	if err != nil {
//		d.logger.Error(op, err)
//		return fmt.Errorf("%s: failed to connect to mongoDB: %w", op, err)
//	}
//
//	d.logger.Infof("%s: success", op)
//	return nil
//}
//
//func (r *repository) Delete(ctx context.Context, id string) error {
//	const op = "coat.storage.Delete"
//
//	oid, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		d.logger.Error(op, err)
//		return fmt.Errorf("%s: failed to convert hex to objectid: %w", op, err)
//	}
//	filter := bson.M{"_id": oid}
//	_, err = d.collection.DeleteOne(ctx, filter)
//	if err != nil {
//		d.logger.Error(op, err)
//		return fmt.Errorf("%s: failed to find coat by id: %w", op, err)
//	}
//
//	d.logger.Infof("%s: success", op)
//	return nil
//}
