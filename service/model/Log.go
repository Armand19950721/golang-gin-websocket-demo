package model

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Id       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Key      string    `gorm:"type:text;index:idx_key;"` 
	Value    string    `gorm:"type:text;"`
	CreateAt time.Time `gorm:"autoCreateTime"`
}
