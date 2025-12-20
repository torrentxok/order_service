package models

import "errors"

type Payment struct {
	Transaction  string `db:"transaction" json:"transaction"`
	RequestID    string `db:"request_id" json:"request_id"`
	Currency     string `db:"currency" json:"currency"`
	Provider     string `db:"provider" json:"provider"`
	Amount       int    `db:"amount" json:"amount"`
	PaymentDT    int64  `db:"payment_dt" json:"payment_dt"`
	Bank         string `db:"bank" json:"bank"`
	DeliveryCost int    `db:"delivery_cost" json:"delivery_cost"`
	GoodsTotal   int    `db:"goods_total" json:"goods_total"`
	CustomFee    int    `db:"custom_fee" json:"custom_fee"`
}

func (p *Payment) Validate() error {
	if p.Transaction == "" {
		return errors.New("transaction is empty")
	}
	if p.Currency == "" {
		return errors.New("currency is empty")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	return nil
}
