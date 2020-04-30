package model

import ()

type Variant struct {
	ID       int `gorm:"column:id" json:"id"`
	Width    int
	Height   int
	FileType string
	Url      string
	Image    Image `gorm:"association_autoupdate:false;association_autocreate:false"`
}
