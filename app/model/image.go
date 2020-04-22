package model

type Image struct {
	ID        int    `gorm:"column:id" json:"id"`
	Url       string `gorm:"UNIQUE_INDEX" json:"url"`
	Alt       string `json:"alt"`
	Copyright string `json:"copyright"`
	Category  string `gorm:"INDEX" json:"category"`
	Variants  []Variant
}
