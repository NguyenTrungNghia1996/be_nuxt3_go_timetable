package seed

import (
	"context"
	"fmt"

	"go-fiber-api/config"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedAdminUnit ensures a default admin unit exists.
func SeedAdminUnit() {
	collection := config.DB.Collection("units")
	var existing models.Unit
	id := AdminUnitID
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&existing)
	if err != mongo.ErrNoDocuments {
		fmt.Println("✅ Admin unit already exists.")
		return
	}
	// ensure we use the same ID as specified
	unit := models.Unit{
		ID:        id,
		Name:      "default units",
		Logo:      "",
		SubDomain: "unti1",
		Active:    true,
	}
	if _, err := collection.InsertOne(context.TODO(), unit); err != nil {
		fmt.Println("❌ Failed to seed admin unit:", err)
		return
	}
	fmt.Println("🚀 Admin unit seeded successfully")
}
