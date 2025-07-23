package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ServiceAccount represents an API token with associated scopes
// for automation or third-party integrations.
type ServiceAccount struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Password string             `json:"password,omitempty" bson:"password"`
	Active   bool               `json:"active" bson:"active"`
}
