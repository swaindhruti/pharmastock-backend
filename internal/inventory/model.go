package inventory

import "time"

type Inventory struct {
	StockistID int64     `json:"stockist_id"`
	MedicineID int64     `json:"medicine_id"`
	CreatedAt  time.Time `json:"created_at"`
}
