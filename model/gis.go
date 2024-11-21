package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Type        string        `bson:"type" json:"type"`
	Coordinates [][][]float64 `bson:"coordinates" json:"coordinates"`
}

type Region struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Province    string             `bson:"province" json:"province"`
	District    string             `bson:"district" json:"district"`
	SubDistrict string             `bson:"sub_district" json:"sub_district"`
	Village     string             `bson:"village" json:"village"`
	Border      Location           `bson:"border" json:"border"`
}