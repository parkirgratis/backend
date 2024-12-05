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
