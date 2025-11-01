package models

import "time"

type Toko struct {
	ID        	uint 			`gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    	uint			`gorm:"column:user_id;unique" json:"user_id"`
	NamaToko  	string			`gorm:"column:nama_toko;size:255" json:"nama_toko"`
	URLFoto   	string			`gorm:"column:url_foto;size:255" json:"photo,omitempty"`
	CreatedAt 	time.Time		`gorm:"column:created_at" json:"created_at"`
	UpdatedAt	time.Time		`gorm:"column:updated_at" json:"updated_at"`

	User		User			`gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	
}

type TokoRequest struct {
	NamaToko	string			`json:"nama_toko"`
	URLFoto		string			`json:"photo"`
}

type TokoResponse struct {
	NamaToko	string			`json:"nama_toko"`
	URLFoto		string			`json:"photo"`
	UserID		uint			`json:"user_id"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
}