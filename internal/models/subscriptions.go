package models

type Subscription struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"` // формат: MM-YYYY
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionSumRequest struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	From        string  `form:"from"` // формат: MM-YYYY
	To          string  `form:"to"`   // формат: MM-YYYY
}
