package dto

import "github.com/FruitsAI/Orange/internal/models"

// ProjectListResult 项目列表结果
type ProjectListResult struct {
	List     []models.Project `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// CreateProjectRequest 创建/更新项目请求
type CreateProjectRequest struct {
	Name           string  `json:"name" binding:"required"`
	Company        string  `json:"company" binding:"required"`
	TotalAmount    float64 `json:"total_amount" binding:"required"`
	Status         string  `json:"status"`
	Type           string  `json:"type" binding:"required"`
	ContractNumber string  `json:"contract_number"`
	ContractDate   string  `json:"contract_date"`
	PaymentMethod  string  `json:"payment_method"`
	StartDate      string  `json:"start_date" binding:"required"`
	EndDate        string  `json:"end_date" binding:"required"`
	Description    string  `json:"description"`
	UserID         int64   `json:"-"`
}
