package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	RoleID    int       `json:"role_id" gorm:"not null;default:1"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

type Role struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type JWTClaim struct {
	ID     uuid.UUID
	Name   string
	Email  string
	RoleID int
	jwt.RegisteredClaims
}
