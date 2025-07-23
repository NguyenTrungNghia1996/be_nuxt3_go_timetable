package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ServiceAccount represents an API token with associated scopes
// for automation or third-party integrations.
type ServiceAccount struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	TokenHash string             `json:"token_hash" bson:"token_hash"`
	Scopes    []string           `json:"scopes" bson:"scopes"`
	Active    bool               `json:"active" bson:"active"`
}
