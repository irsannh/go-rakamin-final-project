package models

import "time"

type Product struct {
	ID            	uint			`gorm:"primaryKey;autoIncrement" json:"id"`
	NamaProduk    	string			`gorm:"column:nama_produk;size:255" json:"nama_produk"`
	Slug          	string			`gorm:"column:slug;size:255" json:"slug"`
	HargaReseller 	int				`gorm:"column:harga_reseller" json:"harga_reseller"`
	HargaKonsumen 	int				`gorm:"column:harga_konsumen" json:"harga_konsumen"`
	Stok          	int				`gorm:"column:stok" json:"stok"`
	Deskripsi     	string			`gorm:"column:deskripsi;type:text" json:"deskripsi"`
	CreatedAt     	time.Time		`gorm:"column:created_at" json:"created_at"`
	UpdatedAt		time.Time		`gorm:"column:updated_at" json:"updated_at"`
	IDToko			uint			`gorm:"column:id_toko" json:"id_toko"`
	IDCategory		uint			`gorm:"column:id_category" json:"category_id"`
	FotoProduk		[]FotoProduk	`gorm:"foreignKey:IDProduk;references:ID;constraint:OnDelete:CASCADE" json:"photos"`
	Category		Category		`gorm:"foreignKey:IDCategory;references:ID;constraint:OnDelete:CASCADE" json:"category"`
	Toko			Toko			`gorm:"foreignKey:IDToko;references:ID;constraint:OnDelete:CASCADE" json:"toko"`
}

func (Product) TableName() string {
	return "produk"
}

type FotoProduk struct {
	ID				uint			`gorm:"primaryKey;autoIncrement" json:"id"`
	IDProduk		uint			`gorm:"column:id_produk" json:"id_produk"`
	URL				string			`gorm:"column:url;size:255" json:"url"`
	CreatedAt		time.Time		`gorm:"column:created_at" json:"created_at"`
	UpdatedAt		time.Time		`gorm:"column:updated_at" json:"updated_at"`
}

func (FotoProduk) TableName() string {
	return "foto_produk"
}

type SuccessUploadProduct struct {
	ID				uint			`json:"id_produk"`
	NamaProduk		string			`json:"nama_produk"`
	IDCategory		uint			`json:"category_id"`
	HargaKonsumen	int				`json:"harga_konsumen"`
	HargaReseller	int				`json:"harga_reseller"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`		
	IDToko			uint			`json:"id_toko"`
	Photos			[]FotoProduk 	`json:"photos"`
}

type ProductResponse struct {
	ID				uint			`json:"id_produk"`
	NamaProduk		string			`json:"nama_produk"`
	Category		Category		`json:"category"`
	HargaKonsumen	int				`json:"harga_konsumen"`
	HargaReseller	int				`json:"harga_reseller"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`		
	IDToko			uint			`json:"id_toko"`
	NamaToko		string			`json:"nama_toko"`
	Photos			[]FotoProduk 	`json:"photos"`
}