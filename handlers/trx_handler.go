package handlers

import (
	"errors"
	"fmt"
	"go_evermos_rakamin_irsan/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllMyTrxHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized",
			})
		}

		var userID uint
		switch v := userIDInt.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case uint:
			userID = v
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid user ID",
			})
		}
		

		var trxs []models.Trx
		err := db.Select("trx.id","trx.id_user","trx.alamat_pengiriman","trx.kode_invoice","trx.method_bayar","trx.harga_total","trx.created_at").Where("trx.id_user = ?", userID).Preload("User", func(tx *gorm.DB) *gorm.DB {
    		return tx.Select("id","nama","email")
  		}).Preload("Alamat", func(tx *gorm.DB) *gorm.DB {
			return tx.Select(
			"id",
			"id_user",
			"judul_alamat",
			"nama_penerima",
			"notelp",          // <- ini yang benar
			"detail_alamat",
			"created_at",
			"updated_at",
			)
  		}).Preload("DetailTrx", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id","id_trx","id_log","id_toko","kuantitas","harga_total","created_at")
		}).Preload("DetailTrx.LogProduk", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id","nama_produk","slug","harga_konsumen")
		}).Order("trx.created_at DESC").Find(&trxs).Error

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		resp := make([]models.TrxLite, 0, len(trxs))
		for _, t := range trxs {
			item := models.TrxLite{
				ID:          t.ID,
				UserID:      t.IDUser,
				AlamatKirim: t.AlamatPengiriman, // atau t.AlamatKirim kalau field kamu demikian
				HargaTotal:  int64(t.HargaTotal),
				KodeInvoice: t.KodeInvoice,
				MethodBayar: t.MethodBayar,
				CreatedAt:   t.CreatedAt,
				User: models.UserLite{
					ID:    t.User.ID,
					Nama:  t.User.Nama,
					Email: t.User.Email,
				},
				Alamat: models.AlamatLite{
					ID:           t.Alamat.ID,
					UserID:       t.Alamat.UserID,
					JudulAlamat:  t.Alamat.JudulAlamat,
					NamaPenerima: t.Alamat.NamaPenerima,
					NoTelp:       t.Alamat.NoTelp,       // dari kolom "notelp" di DB
					DetailAlamat: t.Alamat.DetailAlamat,
					CreatedAt:    t.Alamat.CreatedAt,
					UpdatedAt:    t.Alamat.UpdatedAt,
				},
			}

			// DetailTrx
			details := make([]models.DetailTrxLite, 0, len(t.DetailTrx))
			for _, d := range t.DetailTrx {
				details = append(details, models.DetailTrxLite{
					ID:         d.ID,
					IdTrx:      d.IDTrx,
					IdLog:      d.IDLogProduk,
					IdToko:     d.IDToko,
					Kuantitas:  d.Kuantitas,
					HargaTotal: int64(d.HargaTotal),
					CreatedAt:  d.CreatedAt,
					LogProduk:	models.LogProdukLite{
						ID:            d.LogProduk.ID,
						NamaProduk:    d.LogProduk.NamaProduk,
						Slug:          d.LogProduk.Slug,
						HargaKonsumen: int64(d.LogProduk.HargaKonsumen),
					},
				})
			}
			item.DetailTrx = details

			resp = append(resp, item)
		}
		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": resp,
		})
	}
}

func PostNewTrxHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user_id from JWT
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized",
			})
		}

		var userID uint
		switch v := userIDInt.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case uint:
			userID = v
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid user ID",
			})
		}

		var input models.TrxRequest

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid Input",
			})
		}

		if input.MethodBayar == "" || input.AlamatPengiriman == 0 || len(input.DetailTrx) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid Input",
			})
		}

		// validasi alamat pengiriman
		var alamat models.Alamat
		if err := db.Where("id_user = ?", userID).Where("id = ?", input.AlamatPengiriman).First(&alamat).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Not your address",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
		}

		var harga_total int
		createdLogs := make([]models.LogProduk, 0, len(input.DetailTrx))
		for _, p := range input.DetailTrx {
			var product models.Product
			if err := db.First(&product, p.ProductID).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": "error",
					"message": fmt.Sprintf("product %d not found", p.ProductID),
				})
			}

			harga_total = harga_total + (product.HargaKonsumen * p.Kuantitas)

			log := models.LogProduk{
				IDProduk: product.ID,
				NamaProduk: product.NamaProduk,
				Slug: product.Slug,
				HargaReseller: product.HargaReseller,
				HargaKonsumen: product.HargaKonsumen,
				Stok: product.Stok,
				Deskripsi: product.Deskripsi,
				IDToko: product.IDToko,
				IDCategory: product.IDCategory,
			}
			db.Create(&log)

			createdLogs = append(createdLogs, log)
		}

		trx := models.Trx{
			IDUser: userID,
			AlamatPengiriman: input.AlamatPengiriman,
			HargaTotal: harga_total,
			KodeInvoice: fmt.Sprintf("INV-%d", time.Now().UnixNano()),
			MethodBayar: input.MethodBayar,
		}

		if err := db.Create(&trx).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to create new trx",
			})
		}

		for i, l := range createdLogs {
			item := input.DetailTrx[i]

			detail := models.DetailTrx{
				IDTrx: trx.ID,
				IDLogProduk: l.ID,
				IDToko: l.IDToko,
				Kuantitas: item.Kuantitas,
				HargaTotal: trx.HargaTotal,
			}

			if err := db.Create(&detail).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error",
					"message": "failed to create detail_trx",
				})
			}
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id": trx.ID,
				"kode_invoice": trx.KodeInvoice,
			},
		})
		
	}
}

func GetMyTrxByIdHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized",
			})
		}

		var userID uint
		switch v := userIDInt.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case uint:
			userID = v
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid user ID",
			})
		}
		
		trxIDParam := c.Params("id")
		if trxIDParam == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Missing trx id",
			})
		}
		trxIDu64, err := strconv.ParseUint(trxIDParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid trx id",
			})
		}
		trxID := uint(trxIDu64)

		var t models.Trx
		err = db.
			Select("trx.id","trx.id_user","trx.alamat_pengiriman","trx.kode_invoice","trx.method_bayar","trx.harga_total","trx.created_at").
			Where("trx.id = ? AND trx.id_user = ?", trxID, userID).
			Preload("User", func(tx *gorm.DB) *gorm.DB {
				return tx.Select("id","nama","email")
			}).
			Preload("Alamat", func(tx *gorm.DB) *gorm.DB {
				return tx.Select(
					"id",
					"id_user",
					"judul_alamat",
					"nama_penerima",
					"notelp",
					"detail_alamat",
					"created_at",
					"updated_at",
				)
			}).
			Preload("DetailTrx", func(tx *gorm.DB) *gorm.DB {
				return tx.Select("id","id_trx","id_log","id_toko","kuantitas","harga_total","created_at")
			}).
			Preload("DetailTrx.LogProduk", func(tx *gorm.DB) *gorm.DB {
				return tx.Select("id","nama_produk","slug","harga_konsumen")
			}).
			First(&t).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Transaction not found or not belongs to you",
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		item := models.TrxLite{
			ID:          t.ID,
			UserID:      t.IDUser,
			AlamatKirim: t.AlamatPengiriman,
			HargaTotal:  int64(t.HargaTotal),
			KodeInvoice: t.KodeInvoice,
			MethodBayar: t.MethodBayar,
			CreatedAt:   t.CreatedAt,
			User: models.UserLite{
				ID:    t.User.ID,
				Nama:  t.User.Nama,
				Email: t.User.Email,
			},
			Alamat: models.AlamatLite{
				ID:           t.Alamat.ID,
				UserID:       t.Alamat.UserID,
				JudulAlamat:  t.Alamat.JudulAlamat,
				NamaPenerima: t.Alamat.NamaPenerima,
				NoTelp:       t.Alamat.NoTelp, 
				DetailAlamat: t.Alamat.DetailAlamat,
				CreatedAt:    t.Alamat.CreatedAt,
				UpdatedAt:    t.Alamat.UpdatedAt,
			},
		}

		details := make([]models.DetailTrxLite, 0, len(t.DetailTrx))
		for _, d := range t.DetailTrx {
			details = append(details, models.DetailTrxLite{
				ID:         d.ID,
				IdTrx:      d.IDTrx,
				IdLog:      d.IDLogProduk,
				IdToko:     d.IDToko,
				Kuantitas:  d.Kuantitas,
				HargaTotal: int64(d.HargaTotal),
				CreatedAt:  d.CreatedAt,
				LogProduk: models.LogProdukLite{
					ID:            d.LogProduk.ID,
					NamaProduk:    d.LogProduk.NamaProduk,
					Slug:          d.LogProduk.Slug,
					HargaKonsumen: int64(d.LogProduk.HargaKonsumen),
				},
			})
		}
		item.DetailTrx = details

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   item,
		})
	}
}