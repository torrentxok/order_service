package models

import "errors"

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
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
