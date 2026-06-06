package stockist

type Stockist struct {
	StockistID   int64  `json:"id"`
	OwnerName    string `json:"name" validate:"required, min=3, max=50"`
	BuisnessName string `json:"business_name" validate:"required, min=3, max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,numeric, min=10, max=15"`

	Country string `json:"country" validate:"required"`
	State   string `json:"state" validate:"required"`
	City    string `json:"city" validate:"required"`
	PinCode int    `json:"pin_code" validate:"required,numeric"`
	Address string `json:"address" validate:"required, min=10, max=200"`

	GSTNumber string `json:"gst_number" validate:"required"`
}
