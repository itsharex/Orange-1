package repository

import (
	"github.com/FruitsAI/Orange/internal/database"
	"github.com/FruitsAI/Orange/internal/models"
	"gorm.io/gorm"
)

// ProjectRepository 项目数据仓库
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓库
func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{db: database.GetDB()}
}

// FindByID 根据ID查找项目
func (r *ProjectRepository) FindByID(id int64) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// FindByIDWithPayments 根据ID查找项目（包含收款列表）
func (r *ProjectRepository) FindByIDWithPayments(id int64) (*models.Project, error) {
	var project models.Project
	if err := r.db.Preload("Payments").First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// List 获取项目列表
func (r *ProjectRepository) List(userID int64, status, keyword string, page, pageSize int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{}).Preload("User").Where("user_id = ?", userID)

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR company LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// ListRecent 获取最近项目
func (r *ProjectRepository) ListRecent(userID int64, limit int) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Where("user_id = ?", userID).
		Order("create_time DESC").
		Limit(limit).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// Create 创建项目
func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// Update 更新项目
func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

// Delete 删除项目
func (r *ProjectRepository) Delete(id int64) error {
	return r.db.Delete(&models.Project{}, id).Error
}

// UpdateStatus 更新项目状态
func (r *ProjectRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&models.Project{}).Where("id = ?", id).Update("status", status).Error
}

// GetStats 获取统计数据
func (r *ProjectRepository) GetStats(userID int64) (totalAmount, paidAmount, pendingAmount float64, err error) {
	// 总金额
	r.db.Model(&models.Project{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalAmount)

	// 已收金额
	r.db.Model(&models.Payment{}).
		Joins("JOIN projects ON payments.project_id = projects.id").
		Where("projects.user_id = ? AND payments.status = ?", userID, "paid").
		Select("COALESCE(SUM(payments.amount), 0)").Scan(&paidAmount)

	// 待收金额
	pendingAmount = totalAmount - paidAmount

	return
}

// ExistsByContractNumber 检查合同编号是否存在（限定用户）
func (r *ProjectRepository) ExistsByContractNumber(userID int64, contractNumber string, excludeID int64) (bool, error) {
	var count int64
	query := r.db.Model(&models.Project{}).Where("user_id = ? AND contract_number = ?", userID, contractNumber)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetMaxContractNumberByPrefix 获取指定前缀的最大合同编号（限定用户）
func (r *ProjectRepository) GetMaxContractNumberByPrefix(userID int64, prefix string) (string, error) {
	var contractNumber string
	err := r.db.Model(&models.Project{}).
		Where("user_id = ? AND contract_number LIKE ?", userID, prefix+"%").
		Order("contract_number DESC").
		Limit(1).
		Pluck("contract_number", &contractNumber).Error
	if err != nil {
		return "", err
	}
	return contractNumber, nil
}
