package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username      string     `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Password      string     `json:"-" gorm:"size:100;not null"` // 不返回给前端
	Name          string     `json:"name" gorm:"size:50;not null"`
	Email         string     `json:"email" gorm:"size:100"`
	Phone         string     `json:"phone" gorm:"size:20"`
	Avatar        string     `json:"avatar" gorm:"size:255"`
	Role          string     `json:"role" gorm:"size:20;not null;default:'user'"`
	Department    string     `json:"department" gorm:"size:50"`
	Position      string     `json:"position" gorm:"size:50"`
	Status        int        `json:"status" gorm:"default:1"`
	LastLoginTime *time.Time `json:"last_login_time"`
	CreateTime    time.Time  `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime    time.Time  `json:"update_time" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Project 项目模型
type Project struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string     `json:"name" gorm:"size:100;not null"`
	Company        string     `json:"company" gorm:"size:100;not null"`
	TotalAmount    float64    `json:"total_amount" gorm:"type:real;not null"`
	ReceivedAmount float64    `json:"received_amount" gorm:"type:real;default:0"`
	Status         string     `json:"status" gorm:"size:20;not null"`
	Type           string     `json:"type" gorm:"size:50;not null"`
	ContractNumber string     `json:"contract_number" gorm:"size:50"`
	ContractDate   *time.Time `json:"contract_date" gorm:"type:date"`
	PaymentMethod  string     `json:"payment_method" gorm:"size:30"`
	StartDate      time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate        time.Time  `json:"end_date" gorm:"type:date;not null"`
	Description    string     `json:"description"`
	UserID         int64      `json:"user_id" gorm:"not null;index"`
	CreateTime     time.Time  `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime     time.Time  `json:"update_time" gorm:"autoUpdateTime"`

	// 关联
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:ProjectID"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// Payment 收款模型
type Payment struct {
	ID         int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID  int64      `json:"project_id" gorm:"not null;index"`
	Stage      string     `json:"stage" gorm:"size:50;not null"`
	Amount     float64    `json:"amount" gorm:"type:real;not null"`
	Percentage float64    `json:"percentage" gorm:"type:real"`
	PlanDate   time.Time  `json:"plan_date" gorm:"type:date;not null;index"`
	Status     string     `json:"status" gorm:"size:20;not null;index"`
	ActualDate *time.Time `json:"actual_date" gorm:"type:date"`
	Method     string     `json:"method" gorm:"size:30"`
	Remark     string     `json:"remark" gorm:"size:255"`
	UserID     int64      `json:"user_id" gorm:"not null"`
	CreateTime time.Time  `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time  `json:"update_time" gorm:"autoUpdateTime"`

	// 关联
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}

// Dictionary 字典模型
type Dictionary struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code       string    `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Name       string    `json:"name" gorm:"size:50;not null"`
	Status     int       `json:"status" gorm:"default:1"`
	Remark     string    `json:"remark" gorm:"size:255"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`

	// 关联
	Items []DictionaryItem `json:"items,omitempty" gorm:"foreignKey:DictionaryID"`
}

// TableName 指定表名
func (Dictionary) TableName() string {
	return "dictionaries"
}

// DictionaryItem 字典项模型
type DictionaryItem struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	DictionaryID int64     `json:"dictionary_id" gorm:"not null;index"`
	Label        string    `json:"label" gorm:"size:50;not null"`
	Value        string    `json:"value" gorm:"size:50;not null"`
	Sort         int       `json:"sort" gorm:"default:0"`
	Status       int       `json:"status" gorm:"default:1"`
	Remark       string    `json:"remark" gorm:"size:255"`
	CreateTime   time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (DictionaryItem) TableName() string {
	return "dictionary_item"
}

// Notification 通知模型
type Notification struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      string    `json:"title" gorm:"size:100;not null"`
	Content    string    `json:"content" gorm:"not null"`
	Type       int       `json:"type" gorm:"default:1"` // 1:系统通知, 2:活动通知, 3:私信
	SenderID   int64     `json:"sender_id" gorm:"not null;index"`
	IsGlobal   int       `json:"is_global" gorm:"default:0"` // 0:否, 1:是
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`

	// 非数据库字段，用于前端展示
	IsRead bool  `json:"is_read" gorm:"->"`
	Sender *User `json:"sender,omitempty" gorm:"-"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// UserNotification 用户通知关联表
type UserNotification struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int64      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_notification"`
	NotificationID int64      `json:"notification_id" gorm:"not null;uniqueIndex:idx_user_notification"`
	IsRead         int        `json:"is_read" gorm:"default:0"` // 0:未读, 1:已读
	ReadTime       *time.Time `json:"read_time"`
}

// TableName 指定表名
func (UserNotification) TableName() string {
	return "user_notifications"
}
