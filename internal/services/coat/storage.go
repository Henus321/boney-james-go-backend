package coat

import (
	"context"
	"errors"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

type Storage interface {
	GetAll(ctx context.Context) ([]*Coat, error)
	GetOneByID(ctx context.Context, id string) (*Coat, error)
	Create(ctx context.Context, coat *CreateCoatDTO) error
	Delete(ctx context.Context, id string) error
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) GetAll(ctx context.Context) (coats []*Coat, err error) {
	const op = "coat.storage.GetAll"

	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		d.logger.Error(op, err)
		return coats, fmt.Errorf("%s: failed to find all coats: %w", op, err)
	}

	if err = cursor.All(ctx, &coats); err != nil {
		d.logger.Error(op, err)
		return coats, fmt.Errorf("%s: failed to read all documents from cursor: %w", op, err)
	}

	d.logger.Infof("%s: success", op)
	return coats, nil
}

func (d *db) GetOneByID(ctx context.Context, id string) (coat *Coat, err error) {
	const op = "coat.storage.GetOneByID"

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Error(op, err)
		return coat, fmt.Errorf("%s: failed to convert hex to objectid: %w", op, err)
	}
	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			d.logger.Error(op, err)
			return coat, fmt.Errorf("%s: failed to find coat by id: %w", op, err)
		}

		d.logger.Error(op, err)
		return coat, fmt.Errorf("%s: failed to find coat by id: %w", op, err)
	}

	if err = result.Decode(&coat); err != nil {
		d.logger.Error(op, err)
		return coat, fmt.Errorf("%s: failed to decode coat from DB: %w", op, err)
	}

	d.logger.Infof("%s: success", op)
	return coat, nil
}

func (d *db) Create(ctx context.Context, dto *CreateCoatDTO) error {
	const op = "coat.storage.Create"

	collection := d.collection.Database().Collection("coat")

	_, err := collection.InsertOne(ctx, dto)
	if err != nil {
		d.logger.Error(op, err)
		return fmt.Errorf("%s: failed to connect to mongoDB: %w", op, err)
	}

	d.logger.Infof("%s: success", op)
	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	const op = "coat.storage.Delete"

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		d.logger.Error(op, err)
		return fmt.Errorf("%s: failed to convert hex to objectid: %w", op, err)
	}
	filter := bson.M{"_id": oid}
	_, err = d.collection.DeleteOne(ctx, filter)
	if err != nil {
		d.logger.Error(op, err)
		return fmt.Errorf("%s: failed to find coat by id: %w", op, err)
	}

	d.logger.Infof("%s: success", op)
	return nil
}
