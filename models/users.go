package models

import "time"

type User struct {
	ID           	uint			`gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Nama         	string			`gorm:"column:nama;size:255" json:"nama"`
	KataSandi    	string			`gorm:"column:kata_sandi;size:255" json:"-"`
	NoTelp       	string			`gorm:"column:notelp;size:255;unique" json:"no_telp"`
	TanggalLahir 	time.Time		`gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	JenisKelamin	string			`gorm:"column:jenis_kelamin;size:255" json:"jenis_kelamin"`
	Tentang			string			`gorm:"column:tentang;type:text" json:"tentang"`
	Pekerjaan		string			`gorm:"column:pekerjaan;size:255" json:"pekerjaan"`
	Email			string			`gorm:"column:email;size:255;unique" json:"email"`
	IDProvinsi		string			`gorm:"column:id_provinsi;size:255" json:"id_provinsi"`
	IDKota			string			`gorm:"column:id_kota;size:255" json:"id_kota"`
	IsAdmin			bool			`gorm:"column:is_admin;default:false"`
	CreatedAt		time.Time		`gorm:"column:created_at"`
	UpdatedAt		time.Time		`gorm:"column:updated_at"`
}

type RegisterAndUserRequest struct {
	Nama			string			`json:"nama"`
	KataSandi		string			`json:"kata_sandi"`
	NoTelp			string			`json:"no_telp"`
	TanggalLahir	string			`json:"tanggal_lahir"`
	JenisKelamin	string			`json:"jenis_kelamin"`
	Tentang			string			`json:"tentang"`
	Pekerjaan		string			`json:"pekerjaan"`
	Email			string			`json:"email"`
	IDProvinsi		string			`json:"id_provinsi"`
	IDKota			string			`json:"id_kota"`
}

type LoginRequest struct {
	NoTelp			string			`json:"no_telp"`
	KataSandi		string			`json:"kata_sandi"`
}

type UserResponse struct {
	Nama			string			`json:"nama"`
	NoTelp			string			`json:"no_telp"`
	TanggalLahir	string			`json:"tanggal_lahir"`
	JenisKelamin	string			`json:"jenis_kelamin"`
	Tentang			string			`json:"tentang"`
	Pekerjaan		string			`json:"pekerjaan"`
	Email			string			`json:"email"`
	IDProvinsi		string			`json:"id_provinsi"`
	IDKota			string			`json:"id_kota"`
	CreatedAt	time.Time			`json:"created_at"`
	UpdatedAt	time.Time			`json:"updated_at"`
}
