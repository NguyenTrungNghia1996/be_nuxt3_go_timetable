package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(app *fiber.App, db *mongo.Database) {
	// Initialize repositories and controllers once so they can be reused
	userRepo := repositories.NewUserRepository(db)
	unitRepo := repositories.NewUnitRepository(db)
	saRepo := repositories.NewServiceAccountRepository(db)
	authCtrl := controllers.NewAuthController(userRepo, unitRepo)
	unitCtrl := controllers.NewUnitController(unitRepo)
	saCtrl := controllers.NewServiceAccountController(saRepo)

	// Public routes do not require authentication
	app.Post("/login", authCtrl.Login)
	app.Get("/api/units/by_subdomain", unitCtrl.GetBySubDomain)

	// Protected API group requires JWT authentication
	api := app.Group("/api", middleware.Protected())
	api.Put("/presigned_url", controllers.GetUploadUrl)
	api.Delete("/image", controllers.DeleteImage)

	units := api.Group("/units")
	units.Get("", unitCtrl.List)
	units.Post("", unitCtrl.Create)
	units.Put("", unitCtrl.Update)
	units.Delete("", unitCtrl.Delete)

	sas := api.Group("/service_accounts")
	sas.Get("", saCtrl.List)
	sas.Post("", saCtrl.Create)
	sas.Put("", saCtrl.Update)
	sas.Delete("", saCtrl.Delete)
}
