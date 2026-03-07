package models

import "time"

type Payment struct {
	ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID       int        `gorm:"type:int;not null" json:"order_id"`
	Method        string     `gorm:"type:varchar(50);not null" json:"method"` 
	Status        string     `gorm:"type:varchar(50);not null" json:"status"`
	Amount        int        `gorm:"type:bigint;not null" json:"amount"`
	TransactionID string     `gorm:"type:varchar(255)" json:"transaction_id"` 
	PaymentCode   string     `gorm:"type:varchar(255)" json:"payment_code"`   
	PaymentURL    string     `gorm:"type:text" json:"payment_url"`            
	PaidAt        *time.Time `json:"paid_at"`
	ExpiredAt     *time.Time `json:"expired_at"`
	Order         Order      `gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
