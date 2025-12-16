package kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"github.com/torrentxok/order_service/internal/models"
	"go.uber.org/zap"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.OrderService
	logger  *zap.Logger
}

func NewConsumer(reader *kafka.Reader, svc *service.OrderService, logger *zap.Logger) *Consumer {
	return &Consumer{
		reader:  reader,
		service: svc,
		logger:  logger,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	c.logger.Info("kafka consumer started")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.logger.Info("kafka consumer stopped")
				return nil
			}

			c.logger.Error("kafka read error", zap.Error(err))
			continue
		}

		if err := c.handleMessage(ctx, msg); err != nil {
			c.logger.Error(
				"failed to process message",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
			)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg kafka.Message) error {
	// логирование

	var order models.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		c.logger.Warn("failed to unmarshal message", zap.Error(err))
		return err
	}

	if err := order.Validate(); err != nil {
		c.logger.Warn("order validation failed", zap.Error(err))
		return err
	}

	return c.service.CreateOrder(ctx, &order)
}
