package dto

// PaymentRequest 收款请求
type PaymentRequest struct {
	ProjectID  int64   `json:"project_id" binding:"required"`
	Stage      string  `json:"stage" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Percentage float64 `json:"percentage"`
	PlanDate   string  `json:"plan_date" binding:"required"`
	Status     string  `json:"status"`
	Method     string  `json:"method"`
	Remark     string  `json:"remark"`
	UserID     int64   `json:"-"`
}

// ConfirmPaymentRequest 确认收款请求
type ConfirmPaymentRequest struct {
	ActualDate string `json:"actual_date" binding:"required"`
	Method     string `json:"method"`
}
