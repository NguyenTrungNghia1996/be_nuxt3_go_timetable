package repositories

import (
	"context"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Tìm user theo username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Tạo user mới
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	if !user.Active {
		user.Active = true
	}
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Kiểm tra username đã tồn tại
func (r *UserRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAll returns a paginated list of users filtered by username keyword.
// It also returns the total number of matched documents for pagination.
func (r *UserRepository) GetAll(ctx context.Context, search string, page, limit int64) ([]models.User, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter["username"] = bson.M{"$regex": search, "$options": "i"}
	}

	projection := bson.M{"password": 0}
	opts := options.Find().SetProjection(projection)
	if limit > 0 {
		opts.SetLimit(limit).SetSkip((page - 1) * limit)
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update password theo user ID
func (r *UserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"password": hashedPassword}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// Lấy user theo ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetAllByUnit returns users belonging to a unit with optional search and pagination.
func (r *UserRepository) GetAllByUnit(ctx context.Context, unitID primitive.ObjectID, search string, page, limit int64) ([]models.User, int64, error) {
	filter := bson.M{"unit_id": unitID}
	if search != "" {
		filter["username"] = bson.M{"$regex": search, "$options": "i"}
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

	var users []models.User
	for cursor.Next(ctx) {
		var u models.User
		if err := cursor.Decode(&u); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update modifies a user's profile fields by ID.
func (r *UserRepository) Update(ctx context.Context, id string, user *models.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	set := bson.M{"name": user.Name, "url_avatar": user.UrlAvatar, "active": user.Active}
	if user.Password != "" {
		set["password"] = user.Password
	}
	update := bson.M{"$set": set}
	res, err := r.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a user document by ID.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
