package domain

import (
	"time"
)

type Base struct {
	Id        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime"`
}
