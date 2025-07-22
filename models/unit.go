package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Unit represents a site or department in a multi-site setup.
type Unit struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	Code string             `json:"code" bson:"code"`
}

// UnitResponse is used when returning units to clients.
type UnitResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ToResponse converts a Unit to UnitResponse.
func (u Unit) ToResponse() UnitResponse {
	return UnitResponse{
		ID:   u.ID.Hex(),
		Name: u.Name,
		Code: u.Code,
	}
}
