package handlers

import (
	"fmt"
	"go_evermos_rakamin_irsan/models"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ToSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	s = re.ReplaceAllString(s, "")
	reSpaces := regexp.MustCompile(`[\s-]+`)
	s = reSpaces.ReplaceAllString(s, "-")
	return s
}

func PostProductHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
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
				"status": "error",
				"message": "Invalid User ID",
			})
		}

		var toko models.Toko
		if err := db.Where("user_id = ?", userID).First(&toko).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Toko not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		namaProduk := c.FormValue("nama_produk")
		categoryID, _ := strconv.Atoi(c.FormValue("category_id"))
		hargaReseller, _ := strconv.Atoi(c.FormValue("harga_reseller"))
		hargaKonsumen, _ := strconv.Atoi(c.FormValue("harga_konsumen"))
		stok, _ := strconv.Atoi(c.FormValue("stok"))
		deskripsi := c.FormValue("deskripsi")
		slug := ToSlug(namaProduk)

		product := models.Product{
			NamaProduk: namaProduk,
			HargaReseller: hargaReseller,
			HargaKonsumen: hargaKonsumen,
			Stok: stok,
			Deskripsi: deskripsi,
			IDToko: toko.ID,
			IDCategory: uint(categoryID),
			Slug: slug,
		}

		if err := db.Create(&product).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		var uploadedPhotos []models.FotoProduk

		form, err := c.MultipartForm()
		if err == nil && form.File != nil {
			files := form.File["photos"]
			uploadDir := "./uploads/product_photo"
			os.MkdirAll(uploadDir, os.ModePerm)

			for _, file := range files {
				filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
				filepath := filepath.Join(uploadDir, filename)

				if err := c.SaveFile(file, filepath); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": "error",
						"message": "failed to upload photo: " + err.Error(),
					})
				}

				foto := models.FotoProduk{
					IDProduk: product.ID,
					URL: filename,
				}
				if err := db.Create(&foto).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": "error",
						"message": err.Error(),
					})
				}

				uploadedPhotos = append(uploadedPhotos, foto)
			}
		}

		response := models.SuccessUploadProduct{
			ID: product.ID,
			NamaProduk: product.NamaProduk,
			IDCategory: product.IDCategory,
			IDToko: product.IDToko,
			HargaKonsumen: product.HargaKonsumen,
			HargaReseller: product.HargaReseller,
			Photos: uploadedPhotos,
			CreatedAt: product.CreatedAt,
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "success",
			"details": response,
		})
	}
}

func GetAllProductsHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		namaProduk := c.Query("nama_produk")
		categoryID := c.QueryInt("category_id")
		tokoID := c.QueryInt("toko_id")
		minHarga := c.QueryInt("min_harga")
		maxHarga := c.QueryInt("max_harga")
		limit := c.QueryInt("limit", 10)
		page := c.QueryInt("page", 1)
		offset := (page - 1) * limit

		var products []models.Product

		query := db.Preload("Category").Preload("FotoProduk").Preload("Toko", func (db *gorm.DB) *gorm.DB {
			return db.Select("id", "nama_toko")
		})

		if namaProduk != "" {
			query = query.Where("nama_produk LIKE ?", "%"+namaProduk+"%")
		}

		if categoryID != 0 {
			query = query.Where("id_category = ?", categoryID)
		}

		if tokoID != 0 {
			query = query.Where("id_toko = ?", tokoID)
		}

		if minHarga != 0 {
			query = query.Where("harga_konsumen >= ?", minHarga)
		}
		if maxHarga != 0 {
			query = query.Where("harga_konsumen <= ?", maxHarga)
		}

		query = query.Limit(limit).Offset(offset)

		if err := query.Find(&products).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}


		var response []models.ProductResponse
		for _, p := range products {
			res := models.ProductResponse{
				ID: p.ID,
				NamaProduk: p.NamaProduk,
				Category: p.Category,
				HargaKonsumen: p.HargaKonsumen,
				HargaReseller: p.HargaReseller,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
				IDToko: p.IDToko,
				NamaToko: p.Toko.NamaToko, // sekarang nama toko sudah keload
				Photos: p.FotoProduk,
			}
			response = append(response, res)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":   "success",
			"products": response,
		})
	}
}

func GetProductByIdHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		id := c.Params("id")
		var product models.Product

		if err := db.Preload("Category").Preload("FotoProduk").Preload("Toko", func (db *gorm.DB) *gorm.DB {
			return db.Select("id", "nama_toko")
		}).First(&product, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Cannot get detail",
				})
			}
		}

		response := models.ProductResponse{
			ID: product.ID,
			NamaProduk: product.NamaProduk,
			Category: product.Category,
			HargaKonsumen: product.HargaKonsumen,
			HargaReseller: product.HargaReseller,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			IDToko: product.IDToko,
			NamaToko: product.Toko.NamaToko,
			Photos: product.FotoProduk,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": response,
		})
	}
}

func PutProductHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
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
				"status": "error",
				"message": "Invalid User ID",
			})
		}

		id := c.Params("id")
		var product models.Product

		if err := db.Preload("Toko").Preload("FotoProduk").Preload("Category").First(&product, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Product not found",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		if product.Toko.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"message": "You are not allowed to update this product",
			})
		}
		
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"message": "Invalid form data",
			})
		}

		if v := form.Value["nama_produk"]; len(v) > 0 {
			product.NamaProduk = v[0]
			product.Slug = ToSlug(v[0])
		}

		if v := form.Value["category_id"]; len(v) > 0 {
			if cid, err := strconv.Atoi(v[0]); err == nil {
				product.IDCategory = uint(cid)
			}
		}

		if v := form.Value["harga_reseller"]; len(v) > 0 {
			if hr, err := strconv.Atoi(v[0]); err == nil {
				product.HargaReseller = hr
			}
		}

		if v := form.Value["harga_konsumen"]; len(v) > 0 {
			if hk, err := strconv.Atoi(v[0]); err == nil {
				product.HargaKonsumen = hk
			}
		}

		if v := form.Value["stok"]; len(v) > 0 {
			if s, err := strconv.Atoi(v[0]); err == nil {
				product.Stok = s
			}
		}

		if v := form.Value["deskripsi"]; len(v) > 0 {
			product.Deskripsi = v[0]
		}

		var uploadedPhotos []models.FotoProduk

		if files := form.File["photos"]; len(files) > 0 {
			db.Where("id_produk = ?", product.ID).Delete(&models.FotoProduk{})
			
			uploadDir := "./uploads/product_photo"
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to create upload directory",
				})
			}
			for _, file := range files {
				filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
				filepath := filepath.Join(uploadDir, filename)

				if err := c.SaveFile(file, filepath); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": "error",
						"message": "failed to upload photo: " + err.Error(),
					})
				}

				

				foto := models.FotoProduk{
					IDProduk: product.ID,
					URL: "/uploads/product_photo/" + filename,
				}

				if err := db.Create(&foto).Error; err != nil {
					fmt.Println("DB create error: ", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": "error",
						"message": err.Error(),
					})
				}

				uploadedPhotos = append(uploadedPhotos, foto)
			}

			product.FotoProduk = uploadedPhotos
		}

		

		if err := db.Save(&product).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to update product",
			})
		}

		db.Preload("FotoProduk").Preload("Category").Preload("Toko").First(&product, product.ID)
		
		response := models.ProductResponse{
			ID: product.ID,
			NamaProduk: product.NamaProduk,
			Category: product.Category,
			HargaKonsumen: product.HargaKonsumen,
			HargaReseller: product.HargaReseller,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			IDToko: product.IDToko,
			NamaToko: product.Toko.NamaToko,
			Photos: product.FotoProduk,
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"detail": response,
		})
	}
}

func DeleteProductHandler(db *gorm.DB) fiber.Handler {
	return func (c *fiber.Ctx) error {
		userIDInt := c.Locals("user_id")
		if userIDInt == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "error",
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
				"status": "error",
				"message": "Invalid User ID",
			})
		}
		
		id := c.Params("id")
		var product models.Product

		if err := db.Preload("Toko").First(&product, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": "error",
					"message": "Product not found",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": err.Error(),
			})
		}

		if product.Toko.UserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"message": "You are not allowed to delete this product",
			})
		}

		if len(product.FotoProduk) > 0 {
			if err := db.Where("product_id = ?", product.ID).Delete(&models.FotoProduk{}).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error",
					"message": "Failed to delete product photos",
				})
			}
		}

		if err := db.Delete(&product).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"message": "Failed to delete product",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"message": "Product deleted successfully",
		})
	}
}