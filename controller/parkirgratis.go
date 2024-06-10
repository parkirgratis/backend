package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetLokasi(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	kor, err := atdb.GetAllDoc[[]model.Tempat](config.Mongoconn, "tempat", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, kor)
}

func GetMarker(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	mar, err := atdb.GetOneLatestDoc[model.Koordinat](config.Mongoconn, "marker", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, mar)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, mar)
}

func InsertTempat(db *mongo.Database, col string, namaTempat string, lokasi string, fasilitas string, lon float64, lat float64, gambar string) (insertedID primitive.ObjectID, err error) {
	// Membuat map dari data tempat dengan kunci dan nilai yang sesuai
	tempat := bson.M{
		"nama_tempat": namaTempat, // Nama tempat
		"lokasi":      lokasi,     // Lokasi tempat
		"fasilitas":   fasilitas,  // Fasilitas yang tersedia
		"lon":         lon,        // Longitude tempat
		"lat":         lat,        // Latitude tempat
		"gambar":      gambar,     // URL gambar tempat
	}

	// Menyisipkan data tempat ke dalam koleksi MongoDB yang ditentukan
	result, err := db.Collection(col).InsertOne(context.Background(), tempat)
	if err != nil {
		// Jika terjadi error saat penyisipan, tampilkan error dan hentikan fungsi
		fmt.Printf("InsertTempat: %v\n", err)
		return
	}

	// Mengambil ID dari dokumen yang baru disisipkan
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil // Mengembalikan ID yang disisipkan dan error nil (tidak ada error)
}
func InsertKoordinat(db *mongo.Database, col string, markers [][]float64) (insertedID primitive.ObjectID, err error) {
	// Membuat map dari data koordinat dengan kunci dan nilai yang sesuai
	koordinat := bson.M{
		"markers": markers, // Array dari koordinat
	}

	// Menyisipkan data koordinat ke dalam koleksi MongoDB yang ditentukan
	result, err := db.Collection(col).InsertOne(context.Background(), koordinat)
	if err != nil {
		// Jika terjadi error saat penyisipan, tampilkan error dan hentikan fungsi
		fmt.Printf("InsertKoordinat: %v\n", err)
		return
	}

	// Mengambil ID dari dokumen yang baru disisipkan
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil // Mengembalikan ID yang disisipkan dan error nil (tidak ada error)
}

// PostTempatParkir adalah fungsi yang menangani permintaan POST untuk menyimpan data tempat parkir baru.
func PostTempatParkir(respw http.ResponseWriter, req *http.Request) {
	// Membaca data dari body permintaan
	var data struct {
		NamaTempat string  `json:"nama_tempat"`
		Lokasi     string  `json:"lokasi"`
		Fasilitas  string  `json:"fasilitas"`
		Lon        float64 `json:"lon"`
		Lat        float64 `json:"lat"`
		Gambar     string  `json:"gambar"`
	}
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		// Jika terjadi kesalahan dalam mendekode data, kirimkan pesan kesalahan
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	// Memanggil fungsi InsertTempat untuk menyisipkan data ke dalam database
	insertedID, err := InsertTempat(config.Mongoconn, "tempat_parkir", data.NamaTempat, data.Lokasi, data.Fasilitas, data.Lon, data.Lat, data.Gambar)
	if err != nil {
		// Jika terjadi kesalahan saat menyisipkan data, kirimkan pesan kesalahan
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	// Mengirimkan respons sukses dengan ID dari data yang baru disisipkan
	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Tempat parkir berhasil disimpan dengan ID: %s", insertedID.Hex())})
}

// PostKoordinat adalah fungsi yang menangani permintaan POST untuk menyimpan data koordinat baru.
func PostKoordinat(respw http.ResponseWriter, req *http.Request) {
	// Membaca data dari body permintaan
	var data struct {
		Markers [][]float64 `json:"markers"` // Array dari koordinat
	}
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		// Jika terjadi kesalahan dalam mendekode data, kirimkan pesan kesalahan
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	// Memanggil fungsi InsertKoordinat untuk menyisipkan data ke dalam database
	insertedID, err := InsertKoordinat(config.Mongoconn, "koordinat", data.Markers)
	if err != nil {
		// Jika terjadi kesalahan saat menyisipkan data, kirimkan pesan kesalahan
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	// Mengirimkan respons sukses dengan ID dari data yang baru disisipkan
	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Koordinat berhasil disimpan dengan ID: %s", insertedID.Hex())})
}
