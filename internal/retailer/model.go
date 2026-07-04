package retailer

import "time"

type Retailer struct {
	RetailerID   int64
	OwnerName    string
	BusinessName string
	Email        string
	Phone        string
	Country      string
	State        string
	City         string
	PinCode      string
	Address      string
	GSTNumber    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
