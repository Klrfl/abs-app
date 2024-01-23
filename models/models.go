package models

import (
	"time"

	"github.com/google/uuid"
)

type Drink struct {
	Id         uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Drink_name string    `json:"drink_name"`
	Drink_type string    `json:"drink_type"`
	Cold_price *uint     `json:"cold_price"`
	Hot_price  *uint     `json:"hot_price"`
	Created_at time.Time `json:"created_at" gorm:"default:now()"`
	Updated_at time.Time `json:"updated_at" gorm:"default:now()"`
}

type Member struct {
	Id          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Member_name string    `json:"member_name"`
	Created_at  time.Time `json:"created_at" gorm:"default:now()"`
	Updated_at  time.Time `json:"updated_at" gorm:"default:now()"`
}

type BaseOrder struct {
	Id           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Member_id    uuid.UUID `gorm:"type:uuid"`
	Drink_id     uuid.UUID `gorm:"type:uuid"`
	Created_at   time.Time `gorm:"default:now()"`
	Is_completed bool      `gorm:"default:false"`
	Completed_at time.Time `gorm:"default:null"`
}

type Order struct {
	Id           uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Member_name  string     `json:"member_name"`
	Drink_name   string     `json:"drink_name"`
	Drink_type   *string    `json:"drink_type"`
	Hot_price    *uint      `json:"hot_price"`
	Cold_price   *uint      `json:"cold_price" gorm:"default:now()"`
	Created_at   time.Time  `json:"created_at" gorm:"default:now()"`
	Is_completed bool       `json:"is_completed" gorm:"default:false"`
	Completed_at *time.Time `json:"completed_at" gorm:"default:null"`
}
