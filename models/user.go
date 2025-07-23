// package models
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password,omitempty" bson:"password"`
	Name      string             `json:"name" bson:"name"`
	UrlAvatar string             `json:"url_avatar" bson:"url_avatar"`
	UnitID    primitive.ObjectID `json:"unit_id" bson:"unit_id"`
	Active    bool               `json:"active" bson:"active"`
	IsAdmin   bool               `json:"is_admin" bson:"is_admin"`
}
