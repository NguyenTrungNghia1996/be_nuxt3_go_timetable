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
	authCtrl := controllers.NewAuthController(userRepo)
	unitRepo := repositories.NewUnitRepository(db)
	unitCtrl := controllers.NewUnitController(unitRepo)

	// Public routes do not require authentication
	app.Post("/login", authCtrl.Login)

	// Protected API group requires JWT authentication
	api := app.Group("/api", middleware.Protected())
	api.Put("/presigned_url", controllers.GetUploadUrl)
	api.Delete("/image", controllers.DeleteImage)

	units := api.Group("/units")
	units.Get("", unitCtrl.List)
	units.Post("", unitCtrl.Create)
	units.Get("/by_subdomain", unitCtrl.GetBySubDomain)
	units.Put("", unitCtrl.Update)
	units.Delete("", unitCtrl.Delete)
}
