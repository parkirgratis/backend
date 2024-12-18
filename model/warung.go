package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Warung struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama_Tempat        string             `bson:"nama_tempat,omitempty" json:"nama_tempat,omitempty"`
	Lokasi             string             `bson:"lokasi,omitempty" json:"lokasi,omitempty"`
	Jam_Buka           string             `bson:"jam_buka,omitempty" json:"jam_buka,omitempty"`
	Metode_Pembayaran  []string           `bson:"metode_pembayaran,omitempty" json:"metode_pembayaran,omitempty"`
	Lon         	   float64            `bson:"lon,omitempty" json:"lon,omitempty"`
	Lat         	   float64            `bson:"lat,omitempty" json:"lat,omitempty"`
	Gambar     		   string             `bson:"gambar,omitempty" json:"gambar,omitempty"`
	Province    	   string             `bson:"province,omitempty" json:"province,omitempty"`
    District    	   string             `bson:"district,omitempty" json:"district,omitempty"`
    SubDistrict 	   string             `bson:"sub_district,omitempty" json:"sub_district,omitempty"`
    Village     	   string             `bson:"village,omitempty" json:"village,omitempty"`
}

type KoordinatWarung struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Markers [][]float64 `json:"markers"`
}
