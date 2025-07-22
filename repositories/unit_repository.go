package repositories

import (
	"context"
	"errors"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UnitRepository struct {
	collection *mongo.Collection
}

func NewUnitRepository(db *mongo.Database) *UnitRepository {
	return &UnitRepository{collection: db.Collection("units")}
}

func (r *UnitRepository) Create(ctx context.Context, unit *models.Unit) error {
	unit.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, unit)
	return err
}

func (r *UnitRepository) GetAll(ctx context.Context, search string, page, limit int64) ([]models.Unit, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit).SetSkip((page - 1) * limit)
	}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var units []models.Unit
	for cursor.Next(ctx) {
		var u models.Unit
		if err := cursor.Decode(&u); err != nil {
			return nil, 0, err
		}
		units = append(units, u)
	}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return units, total, nil
}

func (r *UnitRepository) UpdateByID(ctx context.Context, id string, unit *models.Unit) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{
		"name": unit.Name,
		"code": unit.Code,
	}}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("unit not found")
	}
	unit.ID = objID
	return nil
}

func (r *UnitRepository) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("unit not found")
	}
	return nil
}
