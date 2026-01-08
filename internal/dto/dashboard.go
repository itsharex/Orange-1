package dto

// Stats 统计数据
type Stats struct {
	TotalAmount            float64 `json:"total_amount"`
	PaidAmount             float64 `json:"paid_amount"`
	PendingAmount          float64 `json:"pending_amount"`
	OverdueAmount          float64 `json:"overdue_amount"`
	TotalTrend             float64 `json:"total_trend"`
	PaidTrend              float64 `json:"paid_trend"`
	PendingTrend           float64 `json:"pending_trend"`
	OverdueTrend           float64 `json:"overdue_trend"`
	AvgCollectionDays      float64 `json:"avg_collection_days"`
	AvgCollectionDaysTrend float64 `json:"avg_collection_days_trend"`
}

// IncomeTrend 收入趋势
type IncomeTrend struct {
	Labels         []string  `json:"labels"`
	ActualValues   []float64 `json:"actual_values"`
	ExpectedValues []float64 `json:"expected_values"`
}
