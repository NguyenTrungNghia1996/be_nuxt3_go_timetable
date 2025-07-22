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

func SeedDefaultUnit() primitive.ObjectID {
	collection := config.DB.Collection("units")
	var existing models.Unit
	err := collection.FindOne(context.TODO(), bson.M{"code": "default"}).Decode(&existing)
	if err == mongo.ErrNoDocuments {
		id := primitive.NewObjectID()
		unit := models.Unit{
			ID:   id,
			Name: "Default Unit",
			Code: "default",
		}
		if _, err := collection.InsertOne(context.TODO(), unit); err != nil {
			fmt.Println("❌ Failed to seed default unit:", err)
			return primitive.NilObjectID
		}
		fmt.Println("🚀 Default unit seeded")
		return id
	} else if err == nil {
		fmt.Println("✅ Default unit already exists")
		return existing.ID
	} else {
		fmt.Println("❌ Failed checking default unit:", err)
		return primitive.NilObjectID
	}
}
