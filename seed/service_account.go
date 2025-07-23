package seed

import (
	"context"
	"fmt"

	"go-fiber-api/config"
	"go-fiber-api/models"
	"go-fiber-api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedAdminServiceAccount creates an initial admin service account if one doesn't exist.
func SeedAdminServiceAccount() {
	collection := config.DB.Collection("service_accounts")
	var existing models.ServiceAccount
	err := collection.FindOne(context.TODO(), bson.M{"name": "admin"}).Decode(&existing)
	if err != mongo.ErrNoDocuments {
		fmt.Println("✅ Admin service account already exists.")
		return
	}
	password, _ := utils.HashPassword("admin123")
	sa := models.ServiceAccount{
		Name:     "admin",
		Password: password,
		Active:   true,
	}
	if _, err := collection.InsertOne(context.TODO(), sa); err != nil {
		fmt.Println("❌ Failed to seed admin service account:", err)
		return
	}
	fmt.Println("🚀 Admin service account seeded successfully: name=admin password=admin123")
}
