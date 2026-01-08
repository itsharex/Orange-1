package service

import (
	"time"

	"github.com/FruitsAI/Orange/internal/database"
	"github.com/FruitsAI/Orange/internal/dto"
	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/repository"
	"gorm.io/gorm"
)

// ProjectService 项目服务
type ProjectService struct {
	projectRepo *repository.ProjectRepository
	paymentRepo *repository.PaymentRepository
}

// NewProjectService 创建项目服务
func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo: repository.NewProjectRepository(),
		paymentRepo: repository.NewPaymentRepository(),
	}
}

// List 获取项目列表
func (s *ProjectService) List(userID int64, status, keyword string, page, pageSize int) (*dto.ProjectListResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	projects, total, err := s.projectRepo.List(userID, status, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectListResult{
		List:     projects,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Get 获取项目详情
func (s *ProjectService) Get(id int64) (*models.Project, error) {
	return s.projectRepo.FindByIDWithPayments(id)
}

// Create 创建项目
func (s *ProjectService) Create(input dto.CreateProjectRequest) (*models.Project, error) {
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, err
	}
	var contractDate *time.Time
	if input.ContractDate != "" {
		t, err := time.Parse("2006-01-02", input.ContractDate)
		if err != nil {
			return nil, err
		}
		contractDate = &t
	}

	project := &models.Project{
		Name:           input.Name,
		Company:        input.Company,
		TotalAmount:    input.TotalAmount,
		Status:         input.Status,
		Type:           input.Type,
		ContractNumber: input.ContractNumber,
		ContractDate:   contractDate,
		PaymentMethod:  input.PaymentMethod,
		StartDate:      startDate,
		EndDate:        endDate,
		Description:    input.Description,
		UserID:         input.UserID,
	}

	if project.Status == "" {
		project.Status = "active"
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

// Update 更新项目
func (s *ProjectService) Update(id int64, input dto.CreateProjectRequest) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, err
	}
	var contractDate *time.Time
	if input.ContractDate != "" {
		t, err := time.Parse("2006-01-02", input.ContractDate)
		if err != nil {
			return nil, err
		}
		contractDate = &t
	}

	project.Name = input.Name
	project.Company = input.Company
	project.TotalAmount = input.TotalAmount
	project.Status = input.Status
	project.Type = input.Type
	project.ContractNumber = input.ContractNumber
	project.ContractDate = contractDate
	project.PaymentMethod = input.PaymentMethod
	project.StartDate = startDate
	project.EndDate = endDate
	project.Description = input.Description

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

// Delete 删除项目
func (s *ProjectService) Delete(id int64) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. 删除关联款项
		if err := tx.Where("project_id = ?", id).Delete(&models.Payment{}).Error; err != nil {
			return err
		}
		// 2. 删除项目
		if err := tx.Delete(&models.Project{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// Archive 归档项目
func (s *ProjectService) Archive(id int64) error {
	return s.projectRepo.UpdateStatus(id, "archived")
}
