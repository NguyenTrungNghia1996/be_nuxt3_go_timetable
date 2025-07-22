package seed

import (
	"context"
	"fmt"
	"go-fiber-api/config"
	"go-fiber-api/models"
	"go-fiber-api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedAdminUser(unitID primitive.ObjectID) {
	collection := config.DB.Collection("users")
	var existing models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": "admin"}).Decode(&existing)
	if err != mongo.ErrNoDocuments {
		fmt.Println("✅ Admin user already exists.")
		return
	}
	password, _ := utils.HashPassword("admin123")
	groupID, _ := primitive.ObjectIDFromHex("685d01ab5e17ba55d0e349f2")
	admin := models.User{
		Username:   "admin",
		Password:   password,
		Name:       "Administrator",
		UrlAvatar:  "",
		RoleGroups: []primitive.ObjectID{groupID},
		UnitID:     unitID,
	}

	_, err = collection.InsertOne(context.TODO(), admin)
	if err != nil {
		fmt.Println("❌ Failed to seed admin:", err)
		return
	}
	fmt.Println("🚀 Admin user seeded successfully: username=admin password=admin123")
}

// SeedDefaultUser creates a regular user account if not already present.
func SeedDefaultUser(unitID primitive.ObjectID) {
	collection := config.DB.Collection("users")
	var existing models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": "user"}).Decode(&existing)
	if err != mongo.ErrNoDocuments {
		fmt.Println("✅ Regular user already exists.")
		return
	}
	password, _ := utils.HashPassword("user123")
	user := models.User{
		Username:   "user",
		Password:   password,
		Name:       "Default User",
		UrlAvatar:  "",
		RoleGroups: []primitive.ObjectID{},
		UnitID:     unitID,
	}

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println("❌ Failed to seed user:", err)
		return
	}
	fmt.Println("🚀 Regular user seeded successfully: username=user password=user123")
}

// SeedSAUser creates a super admin not tied to any unit.
func SeedSAUser(unitID primitive.ObjectID) {
	collection := config.DB.Collection("users")
	var existing models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": "sa"}).Decode(&existing)
	if err != mongo.ErrNoDocuments {
		fmt.Println("✅ SA user already exists")
		return
	}
	password, _ := utils.HashPassword("sa123")
	groupID, _ := primitive.ObjectIDFromHex("685d01ab5e17ba55d0e349f2")
	sa := models.User{
		Username:   "sa",
		Password:   password,
		Name:       "Super Admin",
		UrlAvatar:  "",
		RoleGroups: []primitive.ObjectID{groupID},
		UnitID:     primitive.NilObjectID,
	}
	if _, err := collection.InsertOne(context.TODO(), sa); err != nil {
		fmt.Println("❌ Failed to seed SA user:", err)
		return
	}
	fmt.Println("🚀 SA user seeded successfully: username=sa password=sa123")
}
