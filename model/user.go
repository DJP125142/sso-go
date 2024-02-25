package model

import "time"

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	EmailVerifiedAt string    `json:"email_verified_at"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:date"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:date"`
}

func (User) TableName() string {
	return "users"
}
