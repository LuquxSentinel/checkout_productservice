package main

import (
	"context"

	"github.com/luqu/productservice/types"
)

type Service interface {
	GetProduct(ctx context.Context, productID string) (*types.Product, error)
	GetAllProducts(ctx context.Context) ([]*types.Product, error)
	// GetOnPromotion(ctx context.Context) ([]*types.Product, error)
	// GetByCategory(ctx context.Context, category string) ([]*types.Product, error)
	// GetBySearch(ctx context.Context, query string) ([]*types.Product, error)
}

type ServiceImpl struct {
	storage Storage
}

func NewService(storage Storage) *ServiceImpl {
	return &ServiceImpl{
		storage: storage,
	}
}

func (s *ServiceImpl) GetProduct(ctx context.Context, productID string) (*types.Product, error) {
	product, err := s.storage.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ServiceImpl) GetAllProducts(ctx context.Context) ([]*types.Product, error) {
	return s.storage.GetAllProducts(ctx)
}
