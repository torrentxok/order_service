package models

import "errors"

type Delivery struct {
	Name    string `db:"name" json:"name"`
	Phone   string `db:"phone" json:"phone"`
	Zip     string `db:"zip" json:"zip"`
	City    string `db:"city" json:"city"`
	Address string `db:"address" json:"address"`
	Region  string `db:"region" json:"region"`
	Email   string `db:"email" json:"email"`
}

func (d *Delivery) Validate() error {
	if d.Name == "" {
		return errors.New("name is empty")
	}
	if d.Phone == "" {
		return errors.New("phone is empty")
	}
	if d.City == "" {
		return errors.New("city is empty")
	}
	if d.Address == "" {
		return errors.New("address is empty")
	}
	return nil
}
