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

func InsertTempat(db *mongo.Database, col string, tempat model.Tempat) (insertedID primitive.ObjectID, err error) {
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

func InsertKoordinat(db *mongo.Database, col string, koordinat model.Koordinat) (insertedID primitive.ObjectID, err error) {
	// Fungsi ini bertujuan untuk menyimpan data koordinat ke dalam koleksi MongoDB yang telah ditentukan.
	// Data koordinat yang akan disimpan adalah array dari koordinat yang terdapat pada objek koordinat.
	data := bson.M{"markers": koordinat.Markers}

	// Melakukan penyisipan data ke dalam koleksi MongoDB.
	result, err := db.Collection(col).InsertOne(context.Background(), data)
	if err != nil {
		// Jika terjadi kesalahan saat menyisipkan data, cetak pesan error dan hentikan eksekusi fungsi.
		fmt.Printf("InsertKoordinat: %v\n", err)
		return
	}

	// Setelah data berhasil disisipkan, ambil ID dari dokumen yang baru saja disisipkan.
	insertedID = result.InsertedID.(primitive.ObjectID)

	// Mengembalikan ID yang telah disisipkan dan tidak ada error yang terjadi (nil).
	return insertedID, nil
}

// PostTempatParkir adalah fungsi yang menangani permintaan POST untuk menyimpan data tempat parkir baru.
func PostTempatParkir(respw http.ResponseWriter, req *http.Request) {
	// Membaca data dari body permintaan
	var data model.Tempat
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		// Jika terjadi kesalahan dalam mendekode data, kirimkan pesan kesalahan
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	// Memanggil fungsi InsertTempat untuk menyisipkan data ke dalam database
	insertedID, err := InsertTempat(config.Mongoconn, "tempat", data)
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
	var data model.Koordinat
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		// Jika terjadi kesalahan dalam mendekode data, kirimkan pesan kesalahan
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}
	// Memanggil fungsi InsertKoordinat untuk menyisipkan data ke dalam database
	insertedID, err := InsertKoordinat(config.Mongoconn, "marker", data)
	if err != nil {
		// Jika terjadi kesalahan saat menyisipkan data, kirimkan pesan kesalahan
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}
	// Mengirimkan respons sukses dengan ID dari data yang baru disisipkan
	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Koordinat berhasil disimpan dengan ID: %s", insertedID.Hex())})
}
