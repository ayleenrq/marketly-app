package models

import "time"

type Seller struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int       `gorm:"type:int;unique;not null" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID" json:"user"`
	StoreName        string    `gorm:"type:varchar(255);not null" json:"store_name"`
	StoreDescription *string   `gorm:"type:text" json:"store_description"`
	IsVerified       bool      `gorm:"type:boolean;default:false" json:"is_verified"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
