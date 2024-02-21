package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID           uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid();onDelete:CASCADE"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;onDelete:CASCADE"`
	User         User           `json:"user"`
	CreatedAt    time.Time      `json:"created_at" gorm:"default:now()"`
	IsCompleted  bool           `json:"is_completed" gorm:"default:false"`
	CompletedAt  time.Time      `json:"completed_at" gorm:"default:null"`
	OrderDetails []*OrderDetail `json:"order_details"`
}

type AnonOrder struct {
	UserName     string         `json:"username"`
	OrderDetails []*OrderDetail `json:"order_details"`
}

type OrderDetail struct {
	OrderID           uuid.UUID `json:"order_id" gorm:"primaryKey;type:uuid;onDelete:CASCADE"`
	MenuID            uuid.UUID `json:"menu_id" gorm:"primaryKey;type:uuid;column:menu_id"`
	MenuName          string    `json:"menu_name" gorm:"column:menu_name"`
	MenuType          string    `json:"menu_type" gorm:"column:menu_type"`
	MenuOptionID      int       `json:"menu_option_id" gorm:"column:menu_option_id"`
	MenuOption        string    `json:"menu_option" gorm:"column:menu_option"`
	MenuOptionValueID int       `json:"menu_option_value_id" gorm:"column:menu_option_value_id"`
	MenuOptionValue   string    `json:"menu_option_value" gorm:"column:menu_option_value"`
	Quantity          int       `json:"quantity" gorm:"column:quantity"`
	TotalPrice        int       `json:"total_price" gorm:"column:total_price"`
}

type BaseOrderDetail struct {
	OrderID           uuid.UUID `json:"order_id" gorm:"primaryKey;type:uuid"`
	MenuID            uuid.UUID `json:"menu_id" gorm:"primaryKey;type:uuid"`
	MenuOptionID      int       `json:"menu_option_id"`
	MenuOptionValueID int       `json:"menu_option_value_id"`
	Quantity          int       `json:"quantity"`
}

func (*BaseOrderDetail) TableName() string {
	return "order_details"
}
