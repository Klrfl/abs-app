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
	Created_at time.Time      `json:"created_at"`
	Updated_at time.Time      `json:"updated_at"`
	Deleted_at gorm.DeletedAt `json:"deleted_at"`
}

type Member struct {
	Id          uuid.UUID      `json:"id" gorm:"primaryKey" gorm:"type:uuid"`
	Member_name string         `json:"member_name"`
	Created_at  time.Time      `json:"created_at"`
	Updated_at  time.Time      `json:"updated_at"`
	Deleted_at  gorm.DeletedAt `json:"deleted_at"`
}

type Order struct {
	Id          uuid.UUID `json:"id" gorm:"primaryKey" gorm:"type:uuid"`
	Member_name string    `json:"member_name"`
	Drink_name  string    `json:"drink_name"`
	Drink_type  string    `json:"drink_type"`
	Hot_price   uint      `json:"hot_price"`
	Cold_price  uint      `json:"cold_price"`
	Created_at  time.Time `json:"created_at"`
}
