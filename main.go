package main

import (
	"go-fiber-api/config"
	"go-fiber-api/routes"
	"go-fiber-api/seed"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file")
		} else {
			log.Println("Loaded .env file")
		}
	}

	// Kết nối MongoDB một lần duy nhất
	config.ConnectDB()

	// Seed default data
	unitID := seed.SeedDefaultUnit()
	seed.SeedRoleGroups(unitID)
	seed.SeedAdminUser(unitID)
	seed.SeedSAUser(primitive.NilObjectID)
	seed.SeedDefaultUser(unitID)
	seed.SeedMenus(unitID)

	app := fiber.New()
	app.Use(cors.New())
	routes.Setup(app, config.DB)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
