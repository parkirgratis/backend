package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Tempat struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama_Tempat string             `bson:"nama_tempat,omitempty" json:"nama_tempat,omitempty"`
	Lokasi      string             `bson:"lokasi,omitempty" json:"lokasi,omitempty"`
	Fasilitas   string             `bson:"fasilitas,omitempty" json:"fasilitas,omitempty"`
	Lon         float64            `bson:"lon,omitempty" json:"lon,omitempty"`
	Lat         float64            `bson:"lat,omitempty" json:"lat,omitempty"`
	Gambar      string            `bson:"gambar,omitempty" json:"gambar,omitempty"`
}

type Koordinat struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Markers [][]float64 `json:"markers"`
}
type Admin struct{
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username,omitempty" json:"username,omitempty"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
}

type Token struct {
	ID			string 				`bson:"_id,omitempty" json:"_id,omitempty"`
	Token 		string 				`bson:"token,omitempty" json:"token,omitempty"`
	AdminID		string				`bson:"admin_id,omitempty" json:"admin_id,omitempty"`
	CreatedAt	time.Time			`bson:"created_at,omitempty" json:"created_at,omitempty"` 
}

type Saran struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Gmail      string             `bson:"gmail" json:"gmail"`
	Nama       string             `bson:"nama" json:"nama"`
	Saran_User string             `bson:"saran_user" json:"saran_user"`
	Tanggal    time.Time          `bson:"tanggal" json:"tanggal"`

	
}

type LoginLog struct {
	ID        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string    `bson:"username,omitempty" json:"username,omitempty"`
	Timestamp time.Time `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	Activity  string    `bson:"activity,omitempry" json:"activity,omitempty"`
}
