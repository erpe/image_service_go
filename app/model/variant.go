package model

import ()

type Variant struct {
	ID       int    `gorm:"column:id" json:"id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Filetype string `json:"filetype"`
	Url      string `json:"url"`
	Filename string `json:"filename"`
	Name     string `json:"name"`
	ImageID  int    `json:"image_id"`
	Image    Image  `gorm:"association_autoupdate:false;association_autocreate:false" json:"image,omitempty"`
}

type PostVariant struct {
	ID        int    `gorm:"column:id" json:"id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Filetype  string `json:"filetype" valid:"required,stringlength(3|5)"`
	Name      string `json:"name" valid:"required,stringlength(3|50)"`
	KeepRatio bool   `gorm:"-" json:"keep_ratio"`
}

func (pv *PostVariant) TableName() string {
	return "variants"
}

type ReadVariant struct {
	ID       int    `gorm:"column:id" json:"id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Filetype string `json:"filetype"`
	Url      string `json:"url"`
	Filename string `json:"filename"`
	Name     string `json:"name"`
	ImageID  int    `json:"image_id"`
}

func (pv *ReadVariant) TableName() string {
	return "variants"
}
