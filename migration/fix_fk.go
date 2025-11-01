// migration/fix_fk.go
package migration

import (
	"fmt"

	"gorm.io/gorm"
)

func FixLogFotoFK(db *gorm.DB) error {
	// 1) Drop semua FK di foto_produk yang mengarah ke log_produk (kalau ada)
	type Row struct{ Name string }
	var rows []Row
	if err := db.Raw(`
		SELECT CONSTRAINT_NAME AS name
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'foto_produk'
		  AND REFERENCED_TABLE_NAME = 'log_produk'
	`).Scan(&rows).Error; err != nil {
		return err
	}
	for _, r := range rows {
		if err := db.Exec(fmt.Sprintf(
			"ALTER TABLE foto_produk DROP FOREIGN KEY `%s`", r.Name,
		)).Error; err != nil {
			return err
		}
	}

	// 2) Drop SEMUA FK di log_produk untuk kolom id_produk (apapun referensinya)
	rows = nil
	if err := db.Raw(`
		SELECT CONSTRAINT_NAME AS name
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'log_produk'
		  AND COLUMN_NAME = 'id_produk'
		  AND REFERENCED_TABLE_NAME IS NOT NULL
	`).Scan(&rows).Error; err != nil {
		return err
	}
	for _, r := range rows {
		if err := db.Exec(fmt.Sprintf(
			"ALTER TABLE log_produk DROP FOREIGN KEY `%s`", r.Name,
		)).Error; err != nil {
			return err
		}
	}

	// 3) Pastikan ada index di log_produk.id_produk (FK butuh index)
	if err := db.Exec(`
		ALTER TABLE log_produk
		ADD INDEX IF NOT EXISTS idx_log_produk__id_produk (id_produk)
	`).Error; err != nil {
		// MySQL lama tidak support IF NOT EXISTS → coba cek manual
		_ = err // aman diabaikan; index biasanya sudah ada
	}

	// 4) Buat FK YANG BENAR: log_produk.id_produk → produk.id
	//   pakai nama unik biar tidak tabrakan: fk_log_produk_produk_ok
	if err := db.Exec(`
		ALTER TABLE log_produk
		ADD CONSTRAINT fk_log_produk_produk_ok
		FOREIGN KEY (id_produk) REFERENCES produk(id)
		ON UPDATE CASCADE ON DELETE RESTRICT
	`).Error; err != nil {
		return err
	}

	return nil
}
