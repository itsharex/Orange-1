package dto

// CreateDictionaryItemRequest 创建字典项请求
type CreateDictionaryItemRequest struct {
	Label string `json:"label" binding:"required"`
	Value string `json:"value" binding:"required"`
	Sort  int    `json:"sort"`
}
