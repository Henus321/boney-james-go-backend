package coat

import "context"

type service struct {
	storage *Storage
}

type Service interface {
	GetAllCoats(ctx context.Context) ([]CoatWithOption, error)
	GetCoatByID(ctx context.Context, id string) (Coat, error)
	CreateCoat(ctx context.Context, dto CreateCoatDTO) error
	DeleteCoat(ctx context.Context, id string) error
}

func NewService(storage *Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetAllCoats(ctx context.Context) ([]CoatWithOption, error) {
	return s.storage.GetAll(ctx)
}

func (s *service) GetCoatByID(ctx context.Context, id string) (Coat, error) {
	return s.storage.GetOneByID(ctx, id)
}

func (s *service) CreateCoat(ctx context.Context, dto CreateCoatDTO) error {
	return s.storage.Create(ctx, dto)
}

func (s *service) DeleteCoat(ctx context.Context, id string) error {
	return s.storage.Delete(ctx, id)
}
