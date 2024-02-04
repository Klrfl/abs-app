package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseMenu struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	TypeID    int       `json:"type_id"`
	Type      MenuType  `json:"type"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

func (BaseMenu) TableName() string {
	return "menu"
}

// full menu, result of join with prices
type Menu struct {
	ID      uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	Iced    int       `json:"iced"`
	Hot     int       `json:"hot"`
	Blend   int       `json:"blend"`
	Regular int       `json:"regular"`
	Plain   int       `json:"plain"`
}

type MenuType struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Type string `json:"type"`
}

type MenuAvailableOption struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Option string `json:"name"`
}

type MenuOptionValue struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	OptionID int    `json:"option_id"`
	Value    string `json:"value"`
}

type Member struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"member_name"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

type Order struct {
	ID           uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MemberID     uuid.UUID      `json:"member_id" gorm:"type:uuid"`
	Member       Member         `json:"member" gorm:""`
	CreatedAt    time.Time      `json:"created_at" gorm:"default:now()"`
	IsCompleted  bool           `json:"is_completed" gorm:"default:false"`
	CompletedAt  time.Time      `json:"completed_at" gorm:"default:null"`
	OrderDetails []*OrderDetail `json:"order_details"`
}

type OrderDetail struct {
	OrderID           uuid.UUID `json:"order_id"`
	MenuID            uuid.UUID `json:"menu_id"`
	MenuName          string    `json:"menu_name"`
	MenuType          string    `json:"menu_type"`
	MenuOptionID      int       `json:"menu_option_id"`
	MenuOption        string    `json:"menu_option"`
	MenuOptionValueID int       `json:"menu_option_value_id"`
	MenuOptionValue   string    `json:"menu_option_value"`
	Quantity          int       `json:"quantity"`
	TotalPrice        int       `json:"total_price"`
}
