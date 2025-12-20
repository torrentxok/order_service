package models

import "errors"

type Item struct {
	ChrtID      int    `db:"chrt_id" json:"chrt_id"`
	TrackNumber string `db:"track_number" json:"track_number"`
	Price       int    `db:"price" json:"price"`
	Rid         string `db:"rid" json:"rid"`
	Name        string `db:"name" json:"name"`
	Sale        int    `db:"sale" json:"sale"`
	Size        string `db:"size" json:"size"`
	TotalPrice  int    `db:"total_price" json:"total_price"`
	NmID        int    `db:"nm_id" json:"nm_id"`
	Brand       string `db:"brand" json:"brand"`
	Status      int    `db:"status" json:"status"`
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
