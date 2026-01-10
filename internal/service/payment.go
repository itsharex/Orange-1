package service

import (
	"time"

	"github.com/FruitsAI/Orange/internal/database"
	"github.com/FruitsAI/Orange/internal/dto"
	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PaymentService 收款服务
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	projectRepo *repository.ProjectRepository
}

// NewPaymentService 创建收款服务
func NewPaymentService() *PaymentService {
	return &PaymentService{
		paymentRepo: repository.NewPaymentRepository(),
		projectRepo: repository.NewProjectRepository(),
	}
}

// ListByProject 根据项目获取收款列表
func (s *PaymentService) ListByProject(projectID int64) ([]models.Payment, error) {
	return s.paymentRepo.ListByProject(projectID)
}

// ListUpcoming 获取即将到期的收款
func (s *PaymentService) ListUpcoming(userID int64, days, limit int) ([]models.Payment, error) {
	return s.paymentRepo.ListUpcoming(userID, days, limit)
}

// ListByDateRange 根据日期范围获取收款
func (s *PaymentService) ListByDateRange(userID int64, startDate, endDate string) ([]models.Payment, error) {
	return s.paymentRepo.ListByDateRange(userID, startDate, endDate)
}

// Create 创建收款
func (s *PaymentService) Create(input dto.PaymentRequest) (*models.Payment, error) {
	planDate, err := time.Parse("2006-01-02", input.PlanDate)
	if err != nil {
		return nil, err
	}

	payment := &models.Payment{
		ProjectID: input.ProjectID,
		Stage:     input.Stage,
		Amount:    input.Amount,
		PlanDate:  planDate,
		Status:    input.Status,
		Method:    input.Method,
		Remark:    input.Remark,
		UserID:    input.UserID,
	}

	if payment.Status == "" {
		payment.Status = "pending"
	}

	// 处理业务规则（日期、百分比）
	if err := s.processPaymentRules(payment); err != nil {
		return nil, err
	}

	// 创建收款记录
	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	// 同步项目总金额
	if err := s.syncProjectReceivedAmount(payment.ProjectID); err != nil {
		return nil, err
	}

	return payment, nil
}

// Update 更新收款
func (s *PaymentService) Update(id int64, input dto.PaymentRequest) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	planDate, err := time.Parse("2006-01-02", input.PlanDate)
	if err != nil {
		return nil, err
	}

	payment.Stage = input.Stage
	payment.Amount = input.Amount
	payment.PlanDate = planDate
	payment.Status = input.Status
	payment.Method = input.Method
	payment.Remark = input.Remark

	// 处理业务规则（日期、百分比）
	if err := s.processPaymentRules(payment); err != nil {
		return nil, err
	}

	// 更新收款记录
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	// 同步项目总金额
	if err := s.syncProjectReceivedAmount(payment.ProjectID); err != nil {
		return nil, err
	}

	return payment, nil
}

// processPaymentRules 处理通用业务规则
func (s *PaymentService) processPaymentRules(payment *models.Payment) error {
	// 1. 处理实际收款日期
	if payment.Status == "paid" && payment.ActualDate == nil {
		// 如果状态为已收款但没有指定实际日期，使用计划日期作为实际日期
		payment.ActualDate = &payment.PlanDate
	}
	if payment.Status != "paid" {
		payment.ActualDate = nil
	}

	// 2. 计算百分比
	project, err := s.projectRepo.FindByID(payment.ProjectID)
	if err != nil {
		return err
	}

	if project.TotalAmount > 0 {
		payment.Percentage = (payment.Amount / project.TotalAmount) * 100
	} else {
		payment.Percentage = 0
	}

	return nil
}

// syncProjectReceivedAmount 同步项目已收款总额
func (s *PaymentService) syncProjectReceivedAmount(projectID int64) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	totalReceived, err := s.paymentRepo.SumPaidByProject(projectID)
	if err != nil {
		return err
	}

	if project.ReceivedAmount != totalReceived {
		project.ReceivedAmount = totalReceived
		if err := s.projectRepo.Update(project); err != nil {
			return err
		}
	}

	return nil
}

// Delete 删除收款
func (s *PaymentService) Delete(id int64) error {
	return s.paymentRepo.Delete(id)
}

// Confirm 确认收款
func (s *PaymentService) Confirm(id int64, actualDate, method string) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. 锁定并获取当前收款记录
		var payment models.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&payment, id).Error; err != nil {
			return err
		}

		// 2. 幂等性检查
		if payment.Status == "paid" {
			return nil
		}

		// 3. 更新收款状态
		if err := tx.Model(&payment).Updates(map[string]interface{}{
			"status":      "paid",
			"actual_date": actualDate,
			"method":      method,
		}).Error; err != nil {
			return err
		}

		// 4. 同步项目总额 (使用 current transaction)
		var totalReceived float64
		if err := tx.Model(&models.Payment{}).
			Where("project_id = ? AND status = ?", payment.ProjectID, "paid").
			Select("COALESCE(SUM(amount), 0)").
			Scan(&totalReceived).Error; err != nil {
			return err
		}

		// 5. 更新项目
		if err := tx.Model(&models.Project{}).
			Where("id = ?", payment.ProjectID).
			Update("received_amount", totalReceived).Error; err != nil {
			return err
		}

		return nil
	})
}
