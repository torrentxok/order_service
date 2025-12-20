package service

import (
	"context"
	"errors"

	"github.com/torrentxok/order_service/internal/cache"
	"github.com/torrentxok/order_service/internal/models"
	"github.com/torrentxok/order_service/internal/repository"
	"go.uber.org/zap"
)

var ErrOrderNotFound = errors.New("order not found")

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

func (s *OrderService) WarmUpCache(ctx context.Context) error {
	orders, err := s.repo.GetLastOrders(ctx, s.cache.Capacity())
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.cache.Set(order.OrderUID, order)
	}

	s.logger.Info("cache warmed up", zap.Int("count", len(orders)))
	return nil
}

// получить запись
func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	if order, ok := s.cache.Get(orderUID); ok {
		s.logger.Debug("order found in cache", zap.String("order_uid", orderUID))
		return order, nil
	}

	order, err := s.repo.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	s.cache.Set(orderUID, order)

	return order, nil
}
