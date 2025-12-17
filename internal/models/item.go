package models

import "errors"

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (i *Item) Validate() error {
	if i.ChrtID <= 0 {
		return errors.New("chrt_id must be positive")
	}
	if i.Name == "" {
		return errors.New("name is empty")
	}
	if i.Price < 0 {
		return errors.New("price must be >= 0")
	}
	return nil
}
