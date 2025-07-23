package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ServiceAccount represents an API token with associated scopes
// for automation or third-party integrations.
type ServiceAccount struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password,omitempty" bson:"password"`
	Name      string             `json:"name" bson:"name"`
	UrlAvatar string             `json:"url_avatar" bson:"url_avatar"`
	Active    bool               `json:"active" bson:"active"`
}
