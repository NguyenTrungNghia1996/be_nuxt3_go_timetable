package repositories

import (
	"context"

	"go-fiber-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceAccountRepository struct {
	collection *mongo.Collection
}

func NewServiceAccountRepository(db *mongo.Database) *ServiceAccountRepository {
	return &ServiceAccountRepository{collection: db.Collection("service_accounts")}
}

func (r *ServiceAccountRepository) Create(ctx context.Context, sa *models.ServiceAccount) error {
	sa.ID = primitive.NewObjectID()
	if !sa.Active {
		sa.Active = true
	}
	_, err := r.collection.InsertOne(ctx, sa)
	return err
}

func (r *ServiceAccountRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ServiceAccountRepository) GetAll(ctx context.Context, search string, page, limit int64) ([]models.ServiceAccount, int64, error) {
	filter := bson.M{}
	if search != "" {
		regex := bson.M{"$regex": search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"username": regex},
			{"name": regex},
		}
	}

	projection := bson.M{"password": 0}
	opts := options.Find().SetProjection(projection)
	if limit > 0 {
		if page <= 0 {
			page = 1
		}
		opts.SetLimit(limit).SetSkip((page - 1) * limit)
	}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var items []models.ServiceAccount
	for cursor.Next(ctx) {
		var sa models.ServiceAccount
		if err := cursor.Decode(&sa); err != nil {
			return nil, 0, err
		}
		items = append(items, sa)
	}
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ServiceAccountRepository) FindByID(ctx context.Context, id string) (*models.ServiceAccount, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var sa models.ServiceAccount
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(&sa); err != nil {
		return nil, err
	}
	return &sa, nil
}

func (r *ServiceAccountRepository) FindByUsername(ctx context.Context, username string) (*models.ServiceAccount, error) {
	var sa models.ServiceAccount
	if err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&sa); err != nil {
		return nil, err
	}
	return &sa, nil
}

func (r *ServiceAccountRepository) Update(ctx context.Context, id string, sa *models.ServiceAccount) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	set := bson.M{"name": sa.Name, "url_avatar": sa.UrlAvatar, "active": sa.Active}
	if sa.Password != "" {
		set["password"] = sa.Password
	}
	update := bson.M{"$set": set}
	_, err = r.collection.UpdateByID(ctx, objID, update)
	return err
}

func (r *ServiceAccountRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
