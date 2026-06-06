package retailer

import "time"

type Retailer struct {
	RetailerID   int    `json:"id"`
	OwnerName    string `json:"name"`
	BuisnessName string `json:"business_name"`
	GSTNumber    string `json:"gst_number"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`

	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	PinCode int    `json:"pin_code"`
	Address string `json:"address"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
