package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Picbot struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name           string    `gorm:"type:text;index:idx_name;"`
	ThemeGroupId   string    `gorm:"type:text;"`
	SocketToken    string    `gorm:"type:text;"`
	SocketTokenOld string    `gorm:"type:text;"`
	Remark         string    `gorm:"type:text;"`
	Used           bool      `gorm:"type:boolean"` // 排序這台picbot是否有用過
	CreateAt       time.Time `gorm:"autoCreateTime"`
	UpdateAt       time.Time `gorm:"autoUpdateTime"`
	DeleteAt       gorm.DeletedAt
}
