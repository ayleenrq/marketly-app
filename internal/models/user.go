package models

import "time"

type User struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID      int       `gorm:"type:int;not null" json:"role_id"`
	Role        Role      `gorm:"foreignKey:RoleID" json:"role"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Email       string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
	PhoneNumber *string   `gorm:"type:varchar(20)" json:"phone_number"`
	PhotoURL    *string   `gorm:"type:text" json:"photo_url"`
	IsActive    bool      `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
