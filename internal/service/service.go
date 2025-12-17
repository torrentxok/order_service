package service

import (
	"context"

	"github.com/torrentxok/order_service/internal/cache"
	"github.com/torrentxok/order_service/internal/models"
	"github.com/torrentxok/order_service/internal/repository"
	"go.uber.org/zap"
)

type OrderService struct {
	repo   repository.OrderRepository
	cache  cache.OrderCache
	logger *zap.Logger
}

func NewOrderService(repo repository.OrderRepository, cache cache.OrderCache, logger *zap.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	exists, err := s.repo.Exists(ctx, order.OrderUID)
	if err != nil {
		return err
	}
	if exists {
		s.logger.Info("order already exists", zap.String("order_uid", order.OrderUID))
		return nil
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return err
	}

	s.cache.Set(order.OrderUID, order)

	return nil
}

// загружаем при запуске в кэш
func (s *OrderService) WarmUpCache(ctx context.Context) error {
	return nil
}

// получить запись
func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	return nil, nil
}
