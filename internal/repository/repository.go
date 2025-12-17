package repository

import (
	"context"

	"github.com/torrentxok/order_service/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
	Exists(ctx context.Context, orderUID string) (bool, error)
	GetLastOrders(ctx context.Context, limit int) ([]*models.Order, error)
}
