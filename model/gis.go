package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Region struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Longitude float64 `bson:"long" json:"long"`
	Latitude  float64 `bson:"lat" json:"lat"`
	Province    string             `bson:"province" json:"province"`
	District    string             `bson:"district" json:"district"`
	SubDistrict string             `bson:"sub_district" json:"sub_district"`
	Village     string             `bson:"village" json:"village"`
	Border      Location           `bson:"border" json:"border"`
}

type Location struct {
	Type        string        `bson:"type" json:"type"`
	Coordinates [][][]float64 `bson:"coordinates" json:"coordinates"`
}

type Geometry struct {
	Type        string       `bson:"type" json:"type"`
	Coordinates [][2]float64 `bson:"coordinates" json:"coordinates"`
}

type LongLat struct {
	Longitude float64 `bson:"long" json:"long"`
	Latitude  float64 `bson:"lat" json:"lat"`
}