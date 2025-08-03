package models

type Subscription struct {
	ID          string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"1299"` // Цена в центах
	UserID      string  `json:"user_id" example:"987e6543-e21b-12d3-a456-426614174999"`
	StartDate   string  `json:"start_date" example:"01-2024"`         // формат: MM-YYYY
	EndDate     *string `json:"end_date,omitempty" example:"12-2024"` // формат: MM-YYYY
}

type SubscriptionSumRequest struct {
	UserID      *string `form:"user_id" example:"987e6543-e21b-12d3-a456-426614174999"`
	ServiceName *string `form:"service_name" example:"Netflix"`
	From        string  `form:"from" example:"01-2024"` // формат: MM-YYYY
	To          string  `form:"to" example:"12-2024"`   // формат: MM-YYYY
}
