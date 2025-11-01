package models

import (
	"time"
)

type Trx struct {
	ID               	uint		`gorm:"primaryKey;autoIncrement" json:"id"`
	IDUser           	uint		`gorm:"column:id_user" json:"user_id"`
	AlamatPengiriman 	uint		`gorm:"column:alamat_pengiriman" json:"alamat_kirim"`
	HargaTotal       	int			`gorm:"column:harga_total" json:"harga_total"`
	KodeInvoice      	string		`gorm:"column:kode_invoice;size:255" json:"kode_invoice"`
	MethodBayar      	string		`gorm:"column:method_bayar;size:255" json:"method_bayar"`
	CreatedAt        	time.Time	`gorm:"column:created_at" json:"created_at"`
	UpdatedAt			time.Time	`gorm:"column:updated_at" json:"updated_at"`
	User				User		`gorm:"foreignKey:IDUser;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Alamat				Alamat		`gorm:"foreignKey:AlamatPengiriman;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"alamat"`
	DetailTrx 			[]DetailTrx `gorm:"foreignKey:IDTrx;references:ID"`
}

func (Trx) TableName() string {
	return "trx"
}


type DetailTrx struct {
	ID					uint		`gorm:"primaryKey;autoIncrement" json:"id"`
	IDTrx				uint		`gorm:"column:id_trx" json:"id_trx"`
	IDLogProduk			uint		`gorm:"column:id_log" json:"id_log"`
	IDToko				uint		`gorm:"column:id_toko" json:"id_toko"`
	Kuantitas			int			`gorm:"column:kuantitas" json:"kuantitas"`
	HargaTotal			int			`gorm:"column:harga_total" json:"harga_total"`
	CreatedAt			time.Time	`gorm:"column:created_at" json:"created_at"`
	UpdatedAt			time.Time	`gorm:"column:updated_at" json:"updated_at"`
	LogProduk			LogProduk	`gorm:"foreignKey:IDLogProduk;references:ID;constraint:OnDelete:CASCADE" json:"log_produk"`
	Toko				Toko		`gorm:"foreignKey:IDToko;references:ID;constraint:OnDelete:CASCADE" json:"toko"`
	Trx					Trx			`gorm:"foreignKey:IDTrx;references:ID;constraint:OnDelete:CASCADE" json:"trx"`
}

func (DetailTrx) TableName() string {
	return "detail_trx"
}

type LogProduk struct {
	ID            	uint			`gorm:"primaryKey;autoIncrement" json:"id"`
	IDProduk		uint			`gorm:"column:id_produk" json:"product_id"`
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
	Category		Category		`gorm:"foreignKey:IDCategory;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Toko			Toko			`gorm:"foreignKey:IDToko;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Product			Product			`gorm:"foreignKey:IDProduk;references:ID" json:"-"`
}

func (LogProduk) TableName() string {
	return "log_produk"
}

type DetailTrxRequest struct {
	ProductID			uint		`json:"product_id"`
	Kuantitas			int			`json:"kuantitas"`
}

type TrxRequest struct {
	MethodBayar      	string				`json:"method_bayar"`
	AlamatPengiriman 	uint				`json:"alamat_kirim"`
	DetailTrx			[]DetailTrxRequest	`json:"detail_trx"`
}

type UserLite struct {
	ID    uint   `json:"id"`
	Nama  string `json:"nama"`
	Email string `json:"email"`
}

type AlamatLite struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	JudulAlamat  string    `json:"judul_alamat"`
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`       
	DetailAlamat string    `json:"detail_alamat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LogProdukLite struct {
	ID            uint   `json:"id"`
	NamaProduk    string `json:"nama_produk"`
	Slug          string `json:"slug"`
	HargaKonsumen int64  `json:"harga_konsumen"`
}

type DetailTrxLite struct {
	ID         uint           `json:"id"`
	IdTrx      uint           `json:"id_trx"`
	IdLog      uint           `json:"id_log"`
	IdToko     uint           `json:"id_toko"`
	Kuantitas  int            `json:"kuantitas"`
	HargaTotal int64          `json:"harga_total"`
	CreatedAt  time.Time      `json:"created_at"`
	LogProduk  LogProdukLite  `json:"log_produk"` 
}

type TrxLite struct {
	ID              uint           `json:"id"`
	UserID          uint           `json:"user_id"`
	AlamatKirim     uint           `json:"alamat_kirim"`     
	HargaTotal      int64          `json:"harga_total"`
	KodeInvoice     string         `json:"kode_invoice"`
	MethodBayar     string         `json:"method_bayar"`
	CreatedAt       time.Time      `json:"created_at"`
	User            UserLite       `json:"user"`
	Alamat          AlamatLite     `json:"alamat"`
	DetailTrx       []DetailTrxLite `json:"detail_trx"`     
}