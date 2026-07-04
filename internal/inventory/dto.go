package inventory

type StockistInfo struct {
	StockistID   int64  `json:"id"`
	OwnerName    string `json:"name"`
	BusinessName string `json:"business_name"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type InventoryResponse struct {
	StockistID int64 `json:"stockist_id"`
	MedicineID int64 `json:"medicine_id"`
}

type StockistsByMedicineResponse struct {
	MedicineID int64          `json:"medicine_id"`
	Stockists  []*StockistInfo `json:"stockists"`
}
