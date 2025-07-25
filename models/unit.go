package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Unit represents a single organization in a multi-tenant SaaS app.
type Unit struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Logo      string             `json:"logo" bson:"logo"`
	SubDomain string             `json:"sub_domain" bson:"sub_domain"`
	Active    bool               `json:"active" bson:"active"`
}
