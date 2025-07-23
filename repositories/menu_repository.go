package repositories

import (
	"context"
	"errors"

	"go-fiber-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuRepository struct {
	collection *mongo.Collection
}

func NewMenuRepository(db *mongo.Database) *MenuRepository {
	return &MenuRepository{collection: db.Collection("menus")}
}

func (r *MenuRepository) Create(ctx context.Context, menu *models.Menu) error {
	menu.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, menu)
	return err
}

// GetAll returns menus optionally filtered by a search keyword.
func (r *MenuRepository) GetAll(ctx context.Context, search string) ([]models.Menu, error) {
	filter := bson.M{}
	if search != "" {
		filter["title"] = bson.M{"$regex": search, "$options": "i"}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []models.Menu
	for cursor.Next(ctx) {
		var m models.Menu
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}
	return menus, nil
}

// GetSAMenus returns menus reserved for service accounts.
func (r *MenuRepository) GetSAMenus(ctx context.Context) ([]models.Menu, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_sa": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []models.Menu
	for cursor.Next(ctx) {
		var m models.Menu
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}
	return menus, nil
}

// GetUnitMenus returns menus for a specific unit.
// Global menus without a unit are also included.
func (r *MenuRepository) GetUnitMenus(ctx context.Context, unitID primitive.ObjectID, isAdmin bool) ([]models.Menu, error) {
	filter := bson.M{
		"is_sa": false,
		"$or": []bson.M{
			{"unit_id": unitID},
			{"unit_id": bson.M{"$exists": false}},
		},
	}
	if !isAdmin {
		filter["is_admin"] = false
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []models.Menu
	for cursor.Next(ctx) {
		var m models.Menu
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}
	return menus, nil
}

func (r *MenuRepository) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("menu not found")
	}
	return nil
}

func (r *MenuRepository) UpdateByID(ctx context.Context, id string, menu *models.Menu) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{
		"title":         menu.Title,
		"key":           menu.Key,
		"url":           menu.URL,
		"icon":          menu.Icon,
		"parent_Id":     menu.ParentID,
		"permissionBit": menu.PermissionBit,
	}}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("menu not found")
	}
	menu.ID = objID
	return nil
}
