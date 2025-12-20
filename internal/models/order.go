package models

import (
	"errors"
	"fmt"
	"time"
)

type Order struct {
	OrderUID        string   `db:"order_uid" json:"order_uid"`
	TrackNumber     string   `db:"track_number" json:"track_number"`
	Entry           string   `db:"entry" json:"entry"`
	Delivery        Delivery `json:"delivery"`
	Payment         Payment  `json:"payment"`
	Items           []Item   `json:"items"`
	Locale          string   `db:"locale" json:"locale"`
	InternalSig     string   `db:"internal_signature" json:"internal_signature"`
	CustomerID      string   `db:"customer_id" json:"customer_id"`
	DeliveryService string   `db:"delivery_service" json:"delivery_service"`
	ShardKey        string   `db:"shardkey" json:"shardkey"`
	SmID            int      `db:"sm_id" json:"sm_id"`
	DateCreated     string   `db:"date_created" json:"date_created"`
	OofShard        string   `db:"oof_shard" json:"oof_shard"`
}

func (o *Order) Validate() error {
	if o == nil {
		return errors.New("order is nil")
	}

	if o.OrderUID == "" {
		return errors.New("order_uid is empty")
	}

	if o.TrackNumber == "" {
		return errors.New("track_number is empty")
	}

	if o.CustomerID == "" {
		return errors.New("customer_id is empty")
	}

	if o.DateCreated == "" {
		return errors.New("date_created is empty")
	}

	parsedTime, err := time.Parse(time.RFC3339, o.DateCreated)
	if err != nil {
		return fmt.Errorf("date_created has invalid format: %w", err)
	}

	if parsedTime.After(time.Now().Add(5 * time.Minute)) {
		return errors.New("date_created is in the future")
	}

	if err := o.Delivery.Validate(); err != nil {
		return fmt.Errorf("delivery validation failed: %w", err)
	}

	if err := o.Payment.Validate(); err != nil {
		return fmt.Errorf("payment validation failed: %w", err)
	}

	if len(o.Items) == 0 {
		return errors.New("items list is empty")
	}

	for i, item := range o.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("items[%d] validation failed: %w", i, err)
		}
	}

	return nil
}
