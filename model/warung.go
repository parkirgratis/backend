package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Warung struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nama_Tempat        string             `bson:"nama_tempat,omitempty" json:"nama_tempat,omitempty"`
	Lokasi             string             `bson:"lokasi,omitempty" json:"lokasi,omitempty"`
	Jam_Buka           string             `bson:"jam_buka,omitempty" json:"jam_buka,omitempty"`
	Metode_Pembayaran  []string           `bson:"metode_pembayaran,omitempty" json:"metode_pembayaran,omitempty"`
	Lon         float64            		  `bson:"lon,omitempty" json:"lon,omitempty"`
	Lat         float64            		  `bson:"lat,omitempty" json:"lat,omitempty"`
	Gambar             string             `bson:"foto_pratinjau,omitempty" json:"foto_pratinjau,omitempty"`
	Tanggal_Verifikasi time.Time          `bson:"tanggal_verifikasi,omitempty" json:"tanggal_verifikasi,omitempty"`
}
