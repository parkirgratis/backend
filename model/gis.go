package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Region struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Border    GeoJSON            `bson:"border" json:"border"`
	Longitude float64            `bson:"longitude" json:"longitude"`
	Latitude  float64            `bson:"latitude" json:"latitude"`
}

// Struct GeoJSON untuk border (polygon/point)
type GeoJSON struct {
	Type        string        `bson:"type" json:"type"`
	Coordinates interface{}   `bson:"coordinates" json:"coordinates"`
}