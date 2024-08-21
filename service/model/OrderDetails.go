package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetails struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrderId      uuid.UUID `gorm:"type:uuid;index:idx_order_id;"`
	QrCodeNumber string    `gorm:"type:text;"`
	PackNameA    string    `gorm:"type:text;"`
	PackNameB    string    `gorm:"type:text;"`
	CreateAt     time.Time `gorm:"autoCreateTime"`
	UpdateAt     time.Time `gorm:"autoUpdateTime"`
	DeleteAt     gorm.DeletedAt
}
