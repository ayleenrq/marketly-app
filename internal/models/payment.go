package models

import "time"

type Payment struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   int        `gorm:"type:int;not null" json:"order_id"`
	Method    string     `gorm:"type:varchar(50)" json:"method"`
	Status    string     `gorm:"type:varchar(50)" json:"status"`
	Amount    int        `gorm:"type:bigint;not null" json:"amount"`
	PaidAt    *time.Time `json:"paid_at"`
	Order     Order      `gorm:"foreignKey:OrderID"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
