package models

import "time"

type OrderItem struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   int       `gorm:"type:int;not null" json:"order_id"`
	ProductID int       `gorm:"type:int;not null" json:"product_id"`
	Quantity  int       `gorm:"type:int;not null" json:"quantity"`
	Price     int       `gorm:"type:bigint;not null" json:"price"`
	Order     Order     `gorm:"foreignKey:OrderID"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
