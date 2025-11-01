package routes

import (
	"go_evermos_rakamin_irsan/handlers"
	"go_evermos_rakamin_irsan/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, jwtSecret string) {
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 		"message": "Hello World!",
	// 	})
	// })

	app.Get("/", handlers.WelcomeHandlers)

	api := app.Group("/api/v1")

	category := api.Group("/category")

	category.Get("/",middleware.JWTMiddleware(jwtSecret), handlers.GetAllCategoriesHandler(db))
	category.Get("/:id", middleware.JWTMiddleware(jwtSecret), handlers.GetCategoryByIdHandler(db))
	category.Post("/", middleware.JWTMiddleware(jwtSecret), middleware.AdminMiddleware(), handlers.PostCategoryHandler(db))
	category.Put("/:id", middleware.JWTMiddleware(jwtSecret), middleware.AdminMiddleware(), handlers.PutCategoryHandler(db))
	category.Delete("/:id", middleware.JWTMiddleware(jwtSecret), middleware.AdminMiddleware(), handlers.DeleteCategoryHandler(db))

	auth := api.Group("/auth")

	auth.Post("/register", handlers.RegisterHandler(db))
	auth.Post("/login", handlers.LoginHandler(db, jwtSecret))

	user := api.Group("/user", middleware.JWTMiddleware(jwtSecret))

	user.Get("/", handlers.GetMyProfileHandler(db))
	user.Put("/", handlers.UpdateMyProfileHandler(db))

	alamat := user.Group("/alamat", middleware.JWTMiddleware(jwtSecret))

	alamat.Get("/", handlers.GetAllMyAlamatHandler(db))
	alamat.Get("/:id", handlers.GetMyAlamatByIdHandler(db))
	alamat.Post("/", handlers.PostAlamatHandler(db))
	alamat.Put("/:id", handlers.UpdateMyAlamatHandler(db))
	alamat.Delete("/:id", handlers.DeleteMyAlamatHandler(db))

	toko := api.Group("/toko", middleware.JWTMiddleware(jwtSecret))

	toko.Get("/", handlers.GetAllTokos(db))
	toko.Get("/my", handlers.GetMyTokoHandler(db))
	toko.Put("/:id_toko", handlers.UpdateMyTokoHandlers(db))
	toko.Get("/:id_toko", handlers.GetTokoByIdHandler(db))

	product := api.Group("/product", middleware.JWTMiddleware(jwtSecret))

	product.Get("/", handlers.GetAllProductsHandler(db))
	product.Get("/:id", handlers.GetProductByIdHandler(db))
	product.Post("/", handlers.PostProductHandler(db))
	product.Put("/:id", handlers.PutProductHandler(db))
	product.Delete("/:id", handlers.DeleteProductHandler(db))

	trx := api.Group("/trx", middleware.JWTMiddleware(jwtSecret))

	trx.Get("/", handlers.GetAllMyTrxHandler(db))
	trx.Get("/:id", handlers.GetMyTrxByIdHandler(db))
	trx.Post("/", handlers.PostNewTrxHandler(db))

	provcity := api.Group("/provcity", middleware.JWTMiddleware(jwtSecret))

	provcity.Get("/listprovincies", handlers.ListProvinciesHandler)
	provcity.Get("/listcities/:prov_id", handlers.ListCitiesByProvinciesHandler)
	provcity.Get("/detailprovince/:prov_id", handlers.ProvinceByIdHandlers)
	provcity.Get("/detailcity/:city_id", handlers.CityByIdHandler)
}