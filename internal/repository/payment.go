package repository

import (
	"time"

	"github.com/FruitsAI/Orange/internal/database"
	"github.com/FruitsAI/Orange/internal/models"
	"gorm.io/gorm"
)

// PaymentRepository 收款数据仓库
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository 创建收款仓库
func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{db: database.GetDB()}
}

// FindByID 根据ID查找收款
func (r *PaymentRepository) FindByID(id int64) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByIDWithProject 根据ID查找收款（包含项目信息）
func (r *PaymentRepository) FindByIDWithProject(id int64) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Preload("Project").First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// ListByProject 根据项目ID获取收款列表
func (r *PaymentRepository) ListByProject(projectID int64) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("project_id = ?", projectID).
		Order("plan_date ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// ListUpcoming 获取即将到期的收款
func (r *PaymentRepository) ListUpcoming(userID int64, days int, limit int) ([]models.Payment, error) {
	var payments []models.Payment
	endDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	if err := r.db.Preload("Project").
		Where("user_id = ? AND status = ? AND plan_date <= ?", userID, "pending", endDate).
		Order("plan_date ASC").
		Limit(limit).
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// ListOverdue 获取逾期收款
func (r *PaymentRepository) ListOverdue(userID int64) ([]models.Payment, error) {
	var payments []models.Payment
	today := time.Now().Format("2006-01-02")

	if err := r.db.Where("user_id = ? AND status = ? AND plan_date < ?", userID, "pending", today).
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// Create 创建收款
func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// Update 更新收款
func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// Delete 删除收款
func (r *PaymentRepository) Delete(id int64) error {
	return r.db.Delete(&models.Payment{}, id).Error
}

// Confirm 确认收款
func (r *PaymentRepository) Confirm(id int64, actualDate, method string) error {
	return r.db.Model(&models.Payment{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "paid",
			"actual_date": actualDate,
			"method":      method,
		}).Error
}

// SumByStatus 按状态统计金额
func (r *PaymentRepository) SumByStatus(userID int64, status string) float64 {
	var sum float64
	r.db.Model(&models.Payment{}).
		Where("user_id = ? AND status = ?", userID, status).
		Select("COALESCE(SUM(amount), 0)").Scan(&sum)
	return sum
}

// SumOverdue 统计逾期金额
func (r *PaymentRepository) SumOverdue(userID int64) float64 {
	var sum float64
	today := time.Now().Format("2006-01-02")
	r.db.Model(&models.Payment{}).
		Where("user_id = ? AND status = ? AND plan_date < ?", userID, "pending", today).
		Select("COALESCE(SUM(amount), 0)").Scan(&sum)
	return sum
}

// ListByDateRange 根据日期范围获取收款列表
func (r *PaymentRepository) ListByDateRange(userID int64, startDate, endDate string) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Preload("Project").
		Where("user_id = ? AND plan_date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("plan_date ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// GetIncomeStats 获取收入统计
func (r *PaymentRepository) GetIncomeStats(userID int64, startDate, endDate, interval string) (map[string]float64, map[string]float64, error) {
	expected := make(map[string]float64)
	actual := make(map[string]float64)

	dateFormat := "%Y-%m-%d"
	if interval == "month" {
		dateFormat = "%Y-%m"
	}

	type Result struct {
		Date  string
		Total float64
	}

	// Expected (plan_date)
	var expectedResults []Result
	if err := r.db.Model(&models.Payment{}).
		Select("strftime('"+dateFormat+"', plan_date) as date, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND plan_date BETWEEN ? AND ?", userID, startDate, endDate).
		Group("date").
		Scan(&expectedResults).Error; err != nil {
		return nil, nil, err
	}
	for _, res := range expectedResults {
		expected[res.Date] = res.Total
	}

	// Actual (actual_date, status = paid)
	var actualResults []Result
	if err := r.db.Model(&models.Payment{}).
		Select("strftime('"+dateFormat+"', actual_date) as date, COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND status = 'paid' AND actual_date BETWEEN ? AND ?", userID, startDate, endDate).
		Group("date").
		Scan(&actualResults).Error; err != nil {
		return nil, nil, err
	}
	for _, res := range actualResults {
		actual[res.Date] = res.Total
	}

	return expected, actual, nil
}

// GetStatsByPeriod 获取指定时间段的统计数据
func (r *PaymentRepository) GetStatsByPeriod(userID int64, startDate, endDate string) (total, paid, overdue, avgPeriod float64, err error) {
	// Total: Expected Revenue (Plan Date in range)
	r.db.Model(&models.Payment{}).
		Where("user_id = ? AND plan_date BETWEEN ? AND ?", userID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&total)

	// Paid: Actual Revenue (Actual Date in range)
	r.db.Model(&models.Payment{}).
		Where("user_id = ? AND status = 'paid' AND actual_date BETWEEN ? AND ?", userID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&paid)

	// Overdue: Plan Date in range, Not Paid, and Plan Date < Now
	now := time.Now().Format("2006-01-02")
	r.db.Model(&models.Payment{}).
		Where("user_id = ? AND status = 'pending' AND plan_date BETWEEN ? AND ? AND plan_date < ?", userID, startDate, endDate, now).
		Select("COALESCE(SUM(amount), 0)").Scan(&overdue)

	// AvgPeriod: Average Collection Days (Actual Date - Plan Date) for items paid in this period
	// Using SQLite julianday
	r.db.Model(&models.Payment{}).
		Where("user_id = ? AND status = 'paid' AND actual_date BETWEEN ? AND ?", userID, startDate, endDate).
		Select("COALESCE(AVG(julianday(actual_date) - julianday(plan_date)), 0)").Scan(&avgPeriod)

	return
}

// SumPaidByProject 计算项目中已支付的总金额
func (r *PaymentRepository) SumPaidByProject(projectID int64) (float64, error) {
	var total float64
	err := r.db.Model(&models.Payment{}).
		Where("project_id = ? AND status = ?", projectID, "paid").
		Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}
