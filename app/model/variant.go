package model

import (
	"time"
)

type Variant struct {
	ID        int `gorm:"column:id" json:"id"`
	Width     int
	Height    int
	FileType  string
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Image     Image
}
