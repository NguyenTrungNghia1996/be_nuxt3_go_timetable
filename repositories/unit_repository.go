package repositories

import (
	"context"
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

func (r *UnitRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Unit, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit).SetSkip((page - 1) * limit)
	}
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var units []models.Unit
	for cursor.Next(ctx) {
		var t models.Unit
		if err := cursor.Decode(&t); err != nil {
			return nil, 0, err
		}
		units = append(units, t)
	}
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return units, total, nil
}

func (r *UnitRepository) FindByID(ctx context.Context, id string) (*models.Unit, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var unit models.Unit
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&unit); err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepository) FindBySubDomain(ctx context.Context, subDomain string) (*models.Unit, error) {
	var unit models.Unit
	err := r.collection.FindOne(ctx, bson.M{"sub_domain": subDomain}).Decode(&unit)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepository) Update(ctx context.Context, id string, unit *models.Unit) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"name": unit.Name, "sub_domain": unit.SubDomain}}
	_, err = r.collection.UpdateByID(ctx, objID, update)
	return err
}

func (r *UnitRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
