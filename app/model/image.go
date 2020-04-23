package model

type Image struct {
	ID        int    `gorm:"column:id" json:"id"`
	Url       string `json:"url"`
	Alt       string `json:"alt"`
	Copyright string `json:"copyright"`
	Category  string `gorm:"INDEX" json:"category"`
	Variants  []Variant
}

/* get a grip on post-data for handler.CreateImage */
type PostImage struct {
	ID        int    `gorm:"column:id" json:"id"`
	Url       string `json:"url"`
	Alt       string `json:"alt"`
	Copyright string `json:"copyright"`
	Category  string `gorm:"INDEX" json:"category"`
	// ignore Data while storing to db
	Data string `gorm:"-" json:"data"`
}

func (pi *PostImage) TableName() string {
	return "images"
}
