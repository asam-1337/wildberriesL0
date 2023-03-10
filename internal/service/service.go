package service

import (
	"context"
	"github.com/asam-1337/wildberriesL0/internal/domain/entity"
)

type OrdersRepository interface {
	SelectById(ctx context.Context, id string) (entity.Order, error)
}

type CacheService interface {
	Load(key string) (value entity.Order, loaded bool)
}

type Service struct {
	repo  OrdersRepository
	cache CacheService
}

func NewService(cache CacheService, repo OrdersRepository) *Service {
	return &Service{
		cache: cache,
		repo:  repo,
	}
}

func (s *Service) GetOrder(ctx context.Context, uid string) (entity.Order, error) {
	order, ok := s.cache.Load(uid)
	if ok {
		return order, nil
	}

	return s.repo.SelectById(ctx, uid)

}
