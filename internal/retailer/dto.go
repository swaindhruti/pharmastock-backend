package retailer

import "time"

type CreateRetailerRequest struct {
	OwnerName    string `json:"name" validate:"required,min=3,max=50"`
	BusinessName string `json:"business_name" validate:"required,min=3,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,numeric,min=10,max=15"`
	Country      string `json:"country" validate:"required"`
	State        string `json:"state" validate:"required"`
	City         string `json:"city" validate:"required"`
	PinCode      int    `json:"pin_code" validate:"required,numeric"`
	Address      string `json:"address" validate:"required,min=10,max=200"`
	GSTNumber    string `json:"gst_number" validate:"required"`
}

func (r *CreateRetailerRequest) ToDomain() *Retailer {
	return &Retailer{
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

type UpdateRetailerRequest struct {
	OwnerName    string `json:"name" validate:"required,min=3,max=50"`
	BusinessName string `json:"business_name" validate:"required,min=3,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,numeric,min=10,max=15"`
	Country      string `json:"country" validate:"required"`
	State        string `json:"state" validate:"required"`
	City         string `json:"city" validate:"required"`
	PinCode      int    `json:"pin_code" validate:"required,numeric"`
	Address      string `json:"address" validate:"required,min=10,max=200"`
	GSTNumber    string `json:"gst_number" validate:"required"`
}

func (r *UpdateRetailerRequest) ToDomain(id int64) *Retailer {
	return &Retailer{
		RetailerID:   id,
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

type RetailerResponse struct {
	ID           int64     `json:"id"`
	OwnerName    string    `json:"name"`
	BusinessName string    `json:"business_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Country      string    `json:"country"`
	State        string    `json:"state"`
	City         string    `json:"city"`
	PinCode      int       `json:"pin_code"`
	Address      string    `json:"address"`
	GSTNumber    string    `json:"gst_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewRetailerResponse(r *Retailer) *RetailerResponse {
	return &RetailerResponse{
		ID:           r.RetailerID,
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
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

type PaginatedRetailerResponse struct {
	Items      []*RetailerResponse `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

func NewPaginatedRetailerResponse(p *PaginatedRetailers) *PaginatedRetailerResponse {
	items := make([]*RetailerResponse, len(p.Items))
	for i, r := range p.Items {
		items[i] = NewRetailerResponse(r)
	}
	return &PaginatedRetailerResponse{
		Items:      items,
		Total:      p.Total,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalPages: p.TotalPages,
	}
}
