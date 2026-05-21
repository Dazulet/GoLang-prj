package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           uint           `gorm:"primaryKey"                          json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null"        json:"username"`
	Email        string         `gorm:"uniqueIndex;size:255;not null"       json:"email"`
	PasswordHash string         `gorm:"not null"                            json:"-"`
	Role         Role           `gorm:"type:varchar(20);default:'user'"     json:"role"`
	Avatar       string         `gorm:"size:500"                            json:"avatar,omitempty"`
	Bio          string         `gorm:"size:1000"                           json:"bio,omitempty"`
	CreatedAt    time.Time      `                                           json:"created_at"`
	UpdatedAt    time.Time      `                                           json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"                               json:"-"`
}
