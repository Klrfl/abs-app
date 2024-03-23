package models

import (
	"github.com/google/uuid"
	"time"
)

type Menu struct {
	ID            uuid.UUID       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string          `json:"name"`
	TypeID        int             `json:"type_id"`
	Type          MenuType        `json:"type"`
	CreatedAt     time.Time       `json:"created_at" gorm:"default:now()"`
	UpdatedAt     time.Time       `json:"updated_at" gorm:"default:now()"`
	VariantValues []*VariantValue `json:"variant_values" gorm:"foreignKey:MenuID;references:ID;constraint:OnDelete:CASCADE"`
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

type InputVariantValue struct {
	OptionID         int `json:"option_id"`
	NewOptionID      int `json:"new_option_id"`
	OptionValueID    int `json:"option_value_id"`
	NewOptionValueID int `json:"new_option_value_id"`
	Price            int `json:"price"`
}

type MenuType struct {
	ID   int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
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
