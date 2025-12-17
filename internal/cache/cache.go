package cache

import "github.com/torrentxok/order_service/internal/models"

type OrderCache interface {
	Get(key string) (*models.Order, bool)
	Set(key string, value *models.Order)
	Delete(key string)
	Capacity() int
}
