package model

import (
	pb "service/protos/PicbotEnum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Orders struct {
	Id             uuid.UUID                  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PicbotName     string                     `gorm:"type:text;index:idx_picbot_name;"`
	PaymentType    pb.PicbotPaymentType       `gorm:"type:int8"`
	ItemPrice      int32                      `gorm:"type:int"`
	TotalPrice     int32                      `gorm:"type:int"`
	AmountPrintOut int32                      `gorm:"type:int"`
	Remark         string                     `gorm:"type:text;"`
	State          pb.PicbotOrderState        `gorm:"type:int8"`
	InvoiceState   pb.PicbotOrderInvoiceState `gorm:"type:int8"`
	CarruerNumber  string                     `gorm:"type:text;"`
	PaymentDate    time.Time
	CreateAt       time.Time `gorm:"autoCreateTime"`
	UpdateAt       time.Time `gorm:"autoUpdateTime"`
	DeleteAt       gorm.DeletedAt
}
