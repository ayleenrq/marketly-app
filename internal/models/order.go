package models

import "time"

type Order struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	BuyerID     int       `gorm:"type:int;not null" json:"buyer_id"`
	TotalAmount int       `gorm:"type:bigint" json:"total_amount"`
	Status      string    `gorm:"type:varchar(50)" json:"status"`
	Buyer       User      `gorm:"foreignKey:BuyerID" json:"buyer"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
