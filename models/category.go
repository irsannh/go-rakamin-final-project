package models

import "time"

type Category struct {
	ID           	uint		`gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NamaCategory 	string		`gorm:"column:nama_category;size:255" json:"nama_category"`
	CreatedAt    	time.Time	`gorm:"column:created_at" json:"created_at"`
	UpdatedAt		time.Time	`gorm:"column:updated_at" json:"updated_at"`
}

type CategoryRequest struct {
	NamaCategory	string		`json:"nama_category"`
}