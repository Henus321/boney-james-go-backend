package shop

import "context"

type service struct {
	storage *Storage
}

type Service interface {
	GetAllShops(ctx context.Context) (*[]ShopWithType, error)
	GetShopByID(ctx context.Context, id string) (*ShopWithType, error)
}

func NewService(storage *Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetAllShops(ctx context.Context) (*[]ShopWithType, error) {
	return s.storage.GetAllShops(ctx)
}

func (s *service) GetShopByID(ctx context.Context, id string) (*ShopWithType, error) {
	return s.storage.GetShopByID(ctx, id)
}
