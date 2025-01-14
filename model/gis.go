package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Region struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Longitude   float64            `bson:"long,omitempty" json:"long,omitempty"`
    Latitude    float64            `bson:"lat,omitempty" json:"lat,omitempty"`
    Province    string             `bson:"province,omitempty" json:"province,omitempty"`
    District    string             `bson:"district,omitempty" json:"district,omitempty"`
    SubDistrict string             `bson:"sub_district,omitempty" json:"sub_district,omitempty"`
    Village     string             `bson:"village,omitempty" json:"village,omitempty"`
}

type LongLat struct {
	Longitude float64 `bson:"long" json:"long"`
	Latitude  float64 `bson:"lat" json:"lat"`
    MaxDistance float64 `bson:"max_distance" json:"max_distance"`
}

type Roads struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Type       string             `bson:"type" json:"type"`
	Geometry   Geometry           `bson:"geometry" json:"geometry"`
	Properties Properties         `bson:"properties" json:"properties"`
}

type Geometry struct {
	Type        string       `bson:"type" json:"type"`
	Coordinates [][2]float64 `bson:"coordinates" json:"coordinates"`
}

type Properties struct {
	OSMID   int64  `bson:"osm_id" json:"osm_id"`
	Name    string `bson:"name" json:"name"`
	Highway string `bson:"highway" json:"highway"`
}

type GeoJSONFitur struct {
    Type       string                 `json:"type"`
    Geometry   map[string]interface{} `json:"geometry"`
    Properties map[string]interface{} `json:"properties"`
}

// GeoJSON untuk koleksi fitur GeoJSON
type GeoJSON struct {
    Type     string           `json:"type"`
    Features []GeoJSONFitur `json:"features"`
}