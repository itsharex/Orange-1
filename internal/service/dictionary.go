package service

import (
	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/repository"
)

// DictionaryService 数据字典服务
// 提供通用字典数据的查询和维护功能，支持字典项的增删改查。
type DictionaryService struct {
	dictRepo *repository.DictionaryRepository
}

// NewDictionaryService 创建字典服务实例
func NewDictionaryService() *DictionaryService {
	return &DictionaryService{
		dictRepo: repository.NewDictionaryRepository(),
	}
}

// List 获取所有字典定义列表
// 返回系统中定义的所有字典类型。
func (s *DictionaryService) List() ([]models.Dictionary, error) {
	return s.dictRepo.List()
}

// GetItems 根据字典编码获取对应的字典项列表
// 用于前端下拉框等选择组件的数据源。
//
// 参数:
//   - code: 字典编码 (如 "project_status")
//
// 返回:
//   - []models.DictionaryItem: 按 sort 排序的字典项列表
func (s *DictionaryService) GetItems(code string) ([]models.DictionaryItem, error) {
	return s.dictRepo.GetItemsByCode(code)
}

// CreateItem 为指定字典创建新选项
//
// 参数:
//   - code: 字典编码 (确定归属哪个字典)
//   - label: 显示名称
//   - value: 数据值
//   - sort: 排序权重 (越小越靠前)
//
// 返回:
//   - *models.DictionaryItem: 创建的字典项
func (s *DictionaryService) CreateItem(code, label, value string, sort int) (*models.DictionaryItem, error) {
	// 1. 查找父级字典
	dict, err := s.dictRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	// 2. 构建实体
	item := &models.DictionaryItem{
		DictionaryID: dict.ID,
		Label:        label,
		Value:        value,
		Sort:         sort,
		Status:       1, // 默认启用
	}

	// 3. 写入数据库
	if err := s.dictRepo.CreateItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateItem 更新字典项信息
//
// 参数:
//   - id: 字典项ID
//   - label: 新的显示名称
//   - value: 新的数据值
//   - sort: 新的排序权重
//
// 返回:
//   - *models.DictionaryItem: 更新后的实体
func (s *DictionaryService) UpdateItem(id int64, label, value string, sort int) (*models.DictionaryItem, error) {
	// 1. 获取现有记录 (确保ID存在且保留DictionaryID等字段)
	item, err := s.dictRepo.FindItemByID(id)
	if err != nil {
		return nil, err
	}

	// 2. 更新字段
	item.Label = label
	item.Value = value
	item.Sort = sort

	// 3. 执行更新
	if err := s.dictRepo.UpdateItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteItem 删除指定字典项
func (s *DictionaryService) DeleteItem(id int64) error {
	return s.dictRepo.DeleteItem(id)
}
