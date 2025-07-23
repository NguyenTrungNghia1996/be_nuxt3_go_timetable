package seed

import (
	"context"
	"fmt"

	"go-fiber-api/config"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
)

// SeedMenus populates the menus collection with default entries using the admin unit.
func SeedMenus() {
	col := config.DB.Collection("menus")
	count, err := col.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("❌ Failed to count menus:", err)
		return
	}
	if count > 0 {
		fmt.Println("✅ Menus already seeded.")
		return
	}

	menus := []interface{}{
		models.Menu{
			Title:         "Dashboard",
			Key:           "dashboard",
			URL:           "/dashboard",
			Icon:          "home",
			PermissionBit: 1,
			UnitID:        &AdminUnitID,
			IsSA:          false,
			IsAdmin:       false,
		},
		models.Menu{
			Title:         "Admin Settings",
			Key:           "admin_settings",
			URL:           "/admin/settings",
			Icon:          "settings",
			PermissionBit: 1,
			UnitID:        &AdminUnitID,
			IsSA:          false,
			IsAdmin:       true,
		},
		models.Menu{
			Title:         "SA Tools",
			Key:           "sa_tools",
			URL:           "/sa/tools",
			Icon:          "tool",
			PermissionBit: 1,
			IsSA:          true,
			IsAdmin:       false,
		},
	}

	if _, err := col.InsertMany(context.TODO(), menus); err != nil {
		fmt.Println("❌ Failed to seed menus:", err)
		return
	}
	fmt.Println("🚀 Menus seeded successfully")
}
