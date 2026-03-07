package models

import "time"

type Order struct {
	ID          int         `gorm:"primaryKey;autoIncrement" json:"id"`
	BuyerID     int         `gorm:"type:int;not null" json:"buyer_id"`
	TotalAmount int         `gorm:"type:bigint" json:"total_amount"`
	Status      string      `gorm:"type:varchar(50)" json:"status"`
	User        User        `gorm:"foreignKey:BuyerID;references:ID" json:"buyer"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Payment     *Payment    `gorm:"foreignKey:OrderID;references:ID" json:"payment"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
