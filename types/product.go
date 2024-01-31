package types

import "time"

type Product struct {
	ID          string
	BrandName   string
	ProductName string
	ExpiryDate  time.Time
	Price       float64
}
