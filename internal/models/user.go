package models

import "time"

type User struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID           int       `gorm:"type:int;not null" json:"role_id"`
	Role             Role      `gorm:"foreignKey:RoleID" json:"role"`
	Username         string    `gorm:"type:varchar(100);unique;not null" json:"username"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	Email            string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password         string    `gorm:"type:varchar(255);not null" json:"-"`
	PhoneNumber      *string   `gorm:"type:varchar(20)" json:"phone_number"`
	Address          *string   `gorm:"type:text" json:"address"`
	PhotoURL         *string   `gorm:"type:text" json:"photo_url"`
	IsActive         bool      `gorm:"type:boolean;default:true" json:"is_active"`
	StoreName        *string   `gorm:"type:varchar(255)" json:"store_name,omitempty"`
	StoreDescription *string   `gorm:"type:text" json:"store_description,omitempty"`
	IsVerified       *bool     `gorm:"default:false" json:"is_verified,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
