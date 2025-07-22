package seed

import (
	"context"
	"fmt"

	"go-fiber-api/config"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var defaultUnitID, _ = primitive.ObjectIDFromHex("687f156bb677e045a4ade130")

func SeedDefaultUnit() primitive.ObjectID {
	collection := config.DB.Collection("units")
	var existing models.Unit
	err := collection.FindOne(context.TODO(), bson.M{"_id": defaultUnitID}).Decode(&existing)
	if err == mongo.ErrNoDocuments {
		unit := models.Unit{
			ID:   defaultUnitID,
			Name: "Default Unit",
			Code: "default",
		}
		if _, err := collection.InsertOne(context.TODO(), unit); err != nil {
			fmt.Println("❌ Failed to seed default unit:", err)
			return primitive.NilObjectID
		}
		fmt.Println("🚀 Default unit seeded")
		return defaultUnitID
	} else if err == nil {
		fmt.Println("✅ Default unit already exists")
		return existing.ID
	} else {
		fmt.Println("❌ Failed checking default unit:", err)
		return primitive.NilObjectID
	}
}
