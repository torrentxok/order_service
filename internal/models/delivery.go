package models

import "errors"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
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
