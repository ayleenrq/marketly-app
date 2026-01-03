package models

import "time"

type Product struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SellerID    int       `gorm:"type:int;not null" json:"seller_id"`
	CategoryID  int       `gorm:"type:int;not null" json:"category_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	Price       int       `gorm:"type:bigint;not null" json:"price"`
	Stock       int       `gorm:"type:int;not null" json:"stock"`
	ImageURL    *string   `gorm:"type:text" json:"image_url"`
	IsActive    bool      `gorm:"type:boolean;default:true" json:"is_active"`
	Seller      Seller    `gorm:"foreignKey:SellerID"`
	Category    Category  `gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
