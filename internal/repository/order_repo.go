package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/torrentxok/order_service/internal/config"
	"github.com/torrentxok/order_service/internal/models"
	"go.uber.org/zap"
)

type OrderRepo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRepository(cfg config.DBConfig, logger *zap.Logger) (*OrderRepo, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("database connected")

	return &OrderRepo{
		db:     db,
		logger: logger,
	}, nil
}

func (r *OrderRepo) CreateOrder(ctx context.Context, o *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	// order
	if err := r.insertOrder(ctx, tx, o); err != nil {
		tx.Rollback()
		return err
	}

	// delivery
	if err := r.insertDelivery(ctx, tx, o.OrderUID, &o.Delivery); err != nil {
		tx.Rollback()
		return err
	}

	// payment
	if err := r.insertPayment(ctx, tx, o.OrderUID, &o.Payment); err != nil {
		tx.Rollback()
		return err
	}

	// items
	if err := r.insertItems(ctx, tx, o.OrderUID, o.Items); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}
	return nil
}

func (r *OrderRepo) insertOrder(ctx context.Context, tx *sql.Tx, o *models.Order) error {
	query := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale,
			internal_signature, customer_id, delivery_service,
			shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := tx.ExecContext(ctx, query,
		o.OrderUID,
		o.TrackNumber,
		o.Entry,
		o.Locale,
		o.InternalSig,
		o.CustomerID,
		o.DeliveryService,
		o.ShardKey,
		o.SmID,
		o.DateCreated,
		o.OofShard,
	)

	if err != nil {
		r.logger.Error("insertOrder failed", zap.Error(err))
		return err
	}
	return nil
}

func (r *OrderRepo) insertDelivery(ctx context.Context, tx *sql.Tx, orderUID string, d *models.Delivery) error {
	query := `
		INSERT INTO delivery (
			order_uid, name, phone, zip,
			city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.ExecContext(ctx, query,
		orderUID,
		d.Name,
		d.Phone,
		d.Zip,
		d.City,
		d.Address,
		d.Region,
		d.Email,
	)

	if err != nil {
		r.logger.Error("insertDelivery failed", zap.Error(err))
		return err
	}
	return nil
}

func (r *OrderRepo) insertPayment(ctx context.Context, tx *sql.Tx, orderUID string, p *models.Payment) error {
	query := `
		INSERT INTO payment (
			order_uid, transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := tx.ExecContext(ctx, query,
		orderUID,
		p.Transaction,
		p.RequestID,
		p.Currency,
		p.Provider,
		p.Amount,
		p.PaymentDT,
		p.Bank,
		p.DeliveryCost,
		p.GoodsTotal,
		p.CustomFee,
	)

	if err != nil {
		r.logger.Error("insertPayment failed", zap.Error(err))
		return err
	}

	return nil
}

func (r *OrderRepo) insertItems(ctx context.Context, tx *sql.Tx, orderUID string, items []models.Item) error {
	query := `
		INSERT INTO items (
    		order_uid, chrt_id, track_number, price,
    		rid, name, sale, size, total_price,
    		nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	for _, it := range items {
		_, err := tx.ExecContext(ctx, query,
			orderUID,
			it.ChrtID,
			it.TrackNumber,
			it.Price,
			it.Rid,
			it.Name,
			it.Sale,
			it.Size,
			it.TotalPrice,
			it.NmID,
			it.Brand,
			it.Status,
		)
		if err != nil {
			r.logger.Error("insertItems failed", zap.Error(err))
			return err
		}
	}

	return nil
}

func (r *OrderRepo) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	var order models.Order

	// order
	queryOrder := `
		SELECT	order_uid, track_number, entry, locale,
				internal_signature, customer_id, delivery_service,
				shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`

	err := r.db.GetContext(ctx, &order, queryOrder, orderUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}

		r.logger.Error("failed to fetch order", zap.String("order_uid", orderUID), zap.Error(err))
		return nil, err
	}

	// delivery
	queryDelivery := `
		SELECT
			name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`

	err = r.db.GetContext(ctx, &order.Delivery, queryDelivery, orderUID)
	if err != nil {
		r.logger.Error("failed to fetch delivery", zap.String("order_uid", orderUID), zap.Error(err))
		return nil, err
	}

	// payment
	queryPayment := `
		SELECT	transaction, request_id, currency, provider, amount,
				payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`

	err = r.db.GetContext(ctx, &order.Payment, queryPayment, orderUID)
	if err != nil {
		r.logger.Error("failed to fetch payment", zap.String("order_uid", orderUID), zap.Error(err))
		return nil, err
	}

	// items
	queryItems := `
		SELECT	chrt_id, track_number, price,
    			rid, name, sale, size, total_price,
    			nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`

	var items []models.Item

	err = r.db.SelectContext(ctx, &items, queryItems, orderUID)
	if err != nil {
		r.logger.Error("failed to fetch items", zap.String("order_uid", orderUID), zap.Error(err))
		return nil, err
	}

	order.Items = items

	return &order, nil
}

func (r *OrderRepo) Exists(ctx context.Context, orderUID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM orders
			WHERE order_uid = $1
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, orderUID).Scan(&exists)
	if err != nil {
		r.logger.Error("failed to check order existence",
			zap.Error(err),
			zap.String("order_uid", orderUID),
		)
		return false, err
	}

	return exists, err
}

func (r *OrderRepo) GetLastOrders(ctx context.Context, limit int) ([]*models.Order, error) {
	query := `
		SELECT order_uid
		FROM orders
		ORDER BY date_created DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		r.logger.Error("failed to get last orders uids", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order

	for rows.Next() {
		var orderUID string

		if err = rows.Scan(&orderUID); err != nil {
			return nil, err
		}

		order, err := r.GetOrder(ctx, orderUID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepo) Close() error {
	return r.db.Close()
}
