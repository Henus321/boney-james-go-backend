package coat

import "context"

type service struct {
	storage *Storage
}

type Service interface {
	GetAllCoats(ctx context.Context) (*[]CoatWithOption, error)
	GetCoatByID(ctx context.Context, id string) (*CoatWithOption, error)
	CreateCoat(ctx context.Context, input CoatCreateInput) error
	DeleteCoat(ctx context.Context, id string) error
	CreateCoatOption(ctx context.Context, input CoatOptionCreateInput) error
}

func NewService(storage *Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetAllCoats(ctx context.Context) (*[]CoatWithOption, error) {
	return s.storage.GetAllCoats(ctx)
}

func (s *service) GetCoatByID(ctx context.Context, id string) (*CoatWithOption, error) {
	return s.storage.GetCoatByID(ctx, id)
}

func (s *service) CreateCoat(ctx context.Context, input CoatCreateInput) error {
	return s.storage.CreateCoat(ctx, input)
}

func (s *service) DeleteCoat(ctx context.Context, id string) error {
	return s.storage.DeleteCoat(ctx, id)
}

func (s *service) CreateCoatOption(ctx context.Context, input CoatOptionCreateInput) error {
	return s.storage.CreateCoatOption(ctx, input)
}
