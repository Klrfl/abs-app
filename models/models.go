package models

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	ID            uuid.UUID       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string          `json:"name"`
	TypeID        int             `json:"type_id"`
	Type          MenuType        `json:"type"`
	CreatedAt     time.Time       `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time       `json:"updated_at" gorm:"default:now()"`
	VariantValues []*VariantValue `json:"variant_values" gorm:"foreignKey:MenuID;references:ID"`
}

func (Menu) TableName() string {
	return "menu"
}

type VariantValue struct {
	MenuID        uuid.UUID           `json:"menu_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OptionID      int                 `json:"option_id" gorm:"primaryKey"`
	OptionValueID int                 `json:"option_value_id" gorm:"primaryKey"`
	Option        MenuAvailableOption `json:"option" gorm:"foreignKey:OptionID"`
	OptionValue   MenuOptionValue     `json:"option_value" gorm:"foreignKey:OptionValueID"`
	Price         int                 `json:"price"`
}

type InputMenu struct {
	Name          string `json:"name"`
	TypeID        int    `json:"type_id"`
	OptionID      int    `json:"option_id"`
	OptionValueID int    `json:"option_value_id"`
	Price         int    `json:"price"`
}

func (InputMenu) TableName() string {
	return "menu"
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
	Member       Member         `json:"member"`
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
