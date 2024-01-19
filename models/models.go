package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Drink struct {
	Id         uuid.UUID      `json:"id" gorm:"primaryKey" gorm:"type:uuid"`
	Drink_name string         `json:"drink_name"`
	Drink_type string         `json:"drink_type"`
	Cold_price uint           `json:"cold_price"`
	Hot_price  uint           `json:"hot_price"`
	Created_at time.Time      `json:"created_at" json:"-"`
	Updated_at time.Time      `json:"updated_at" json:"-"`
	Deleted_at gorm.DeletedAt `json:"deleted_at" json:"-"`
}
