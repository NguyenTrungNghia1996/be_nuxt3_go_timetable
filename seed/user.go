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

func SeedAdminUser() {
	collection := config.DB.Collection("users")
	var existing models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": "admin"}).Decode(&existing)
	if err != mongo.ErrNoDocuments {
		fmt.Println("✅ Admin user already exists.")
		return
	}
	password, _ := utils.HashPassword("admin123")
	unitID := AdminUnitID
	admin := models.User{
		Username:  "admin",
		Password:  password,
		Name:      "Administrator",
		UrlAvatar: "",
		Active:    true,
		UnitID:    unitID,
		IsAdmin:   true,
	}

	_, err = collection.InsertOne(context.TODO(), admin)
	if err != nil {
		fmt.Println("❌ Failed to seed admin:", err)
		return
	}
	fmt.Println("🚀 Admin user seeded successfully: username=admin password=admin123")
}

// SeedDefaultUser creates a regular user account if not already present.
// func SeedDefaultUser() {
// 	collection := config.DB.Collection("users")
// 	var existing models.User
// 	err := collection.FindOne(context.TODO(), bson.M{"username": "user"}).Decode(&existing)
// 	if err != mongo.ErrNoDocuments {
// 		fmt.Println("✅ Regular user already exists.")
// 		return
// 	}
// 	password, _ := utils.HashPassword("user123")
//      user := models.User{
//              Username:  "user",
//              Password:  password,
//              Name:      "Default User",
//              UrlAvatar: "",
//      }

// 	_, err = collection.InsertOne(context.TODO(), user)
// 	if err != nil {
// 		fmt.Println("❌ Failed to seed user:", err)
// 		return
// 	}
// 	fmt.Println("🚀 Regular user seeded successfully: username=user password=user123")
// }
