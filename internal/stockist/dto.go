package stockist

import "time"

type CreateStockistRequest struct {
	OwnerName    string `json:"name" validate:"required,min=3,max=50"`
	BusinessName string `json:"business_name" validate:"required,min=3,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,numeric,min=10,max=15"`
	Country      string `json:"country" validate:"required"`
	State        string `json:"state" validate:"required"`
	City         string `json:"city" validate:"required"`
	PinCode      string `json:"pin_code" validate:"required"`
	Address      string `json:"address" validate:"required,min=10,max=200"`
	GSTNumber    string `json:"gst_number" validate:"required"`
}

func (r *CreateStockistRequest) ToDomain() *Stockist {
	return &Stockist{
		OwnerName:    r.OwnerName,
		BusinessName: r.BusinessName,
		Email:        r.Email,
		Phone:        r.Phone,
		Country:      r.Country,
		State:        r.State,
		City:         r.City,
		PinCode:      r.PinCode,
		Address:      r.Address,
		GSTNumber:    r.GSTNumber,
	}
}

type UpdateStockistRequest struct {
	OwnerName    string `json:"name" validate:"required,min=3,max=50"`
	BusinessName string `json:"business_name" validate:"required,min=3,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,numeric,min=10,max=15"`
	Country      string `json:"country" validate:"required"`
	State        string `json:"state" validate:"required"`
	City         string `json:"city" validate:"required"`
	PinCode      string `json:"pin_code" validate:"required"`
	Address      string `json:"address" validate:"required,min=10,max=200"`
	GSTNumber    string `json:"gst_number" validate:"required"`
}

func (r *UpdateStockistRequest) ToDomain(id int64) *Stockist {
	return &Stockist{
		StockistID:   id,
		OwnerName:    r.OwnerName,
		BusinessName: r.BusinessName,
		Email:        r.Email,
		Phone:        r.Phone,
		Country:      r.Country,
		State:        r.State,
		City:         r.City,
		PinCode:      r.PinCode,
		Address:      r.Address,
		GSTNumber:    r.GSTNumber,
	}
}

type StockistResponse struct {
	ID           int64     `json:"id"`
	OwnerName    string    `json:"name"`
	BusinessName string    `json:"business_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Country      string    `json:"country"`
	State        string    `json:"state"`
	City         string    `json:"city"`
	PinCode      string    `json:"pin_code"`
	Address      string    `json:"address"`
	GSTNumber    string    `json:"gst_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewStockistResponse(s *Stockist) *StockistResponse {
	return &StockistResponse{
		ID:           s.StockistID,
		OwnerName:    s.OwnerName,
		BusinessName: s.BusinessName,
		Email:        s.Email,
		Phone:        s.Phone,
		Country:      s.Country,
		State:        s.State,
		City:         s.City,
		PinCode:      s.PinCode,
		Address:      s.Address,
		GSTNumber:    s.GSTNumber,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

type PaginatedStockistResponse struct {
	Items      []*StockistResponse `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

func NewPaginatedStockistResponse(p *PaginatedStockists) *PaginatedStockistResponse {
	items := make([]*StockistResponse, len(p.Items))
	for i, s := range p.Items {
		items[i] = NewStockistResponse(s)
	}
	return &PaginatedStockistResponse{
		Items:      items,
		Total:      p.Total,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalPages: p.TotalPages,
	}
}
