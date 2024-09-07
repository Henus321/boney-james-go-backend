package shop

import "context"

type service struct {
	storage *Storage
}

type Service interface {
	//GetAllCoats(ctx context.Context) (*[]CoatWithOption, error)
	GetShopByID(ctx context.Context, id string) (*ShopWithType, error)
}

func NewService(storage *Storage) Service {
	return &service{storage: storage}
}

//func (s *service) GetAllCoats(ctx context.Context) (*[]CoatWithOption, error) {
//	return s.storage.GetAllCoats(ctx)
//}

func (s *service) GetShopByID(ctx context.Context, id string) (*ShopWithType, error) {
	return s.storage.GetShopByID(ctx, id)
}
