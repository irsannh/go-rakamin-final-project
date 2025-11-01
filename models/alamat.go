package models

import "time"

type Alamat struct {
	ID           	uint   		`gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       	uint   		`gorm:"column:id_user" json:"user_id"`
	JudulAlamat  	string 		`gorm:"column:judul_alamat;size:255" json:"judul_alamat"`
	NamaPenerima 	string 		`gorm:"column:nama_penerima;size:255" json:"nama_penerima"`
	NoTelp       	string 		`gorm:"column:notelp;size:255" json:"no_telp"`
	DetailAlamat 	string 		`gorm:"column:detail_alamat;size:255" json:"detail_alamat"`
	CreatedAt    	time.Time	`gorm:"column:created_at" json:"created_at"`
	UpdatedAt		time.Time	`gorm:"column:updated_at" json:"updated_at"`

	User			User		`gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
}

type AlamatRequest struct {
	JudulAlamat  	string 		`json:"judul_alamat"`
	NamaPenerima 	string 		`json:"nama_penerima"`
	NoTelp       	string 		`json:"no_telp"`
	DetailAlamat 	string 		`json:"detail_alamat"`
}

type AlamatResponse struct {
	JudulAlamat  	string 		`json:"judul_alamat"`
	NamaPenerima 	string 		`json:"nama_penerima"`
	NoTelp       	string 		`json:"no_telp"`
	DetailAlamat 	string 		`json:"detail_alamat"`
}

