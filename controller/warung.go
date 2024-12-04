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
)

func GetaAllWarung(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	warung, err := atdb.GetAllDoc[[]model.Warung](config.Mongoconn, "warung", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, warung)

}

func PostTempatWarung(respw http.ResponseWriter, req *http.Request) {
	var tempatWarung model.Warung

	if err := json.NewDecoder(req.Body).Decode(&tempatWarung); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
	}

	if tempatWarung.Gambar != "" {
		tempatWarung.Gambar = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + tempatWarung.Gambar
	}

	result, err := config.Mongoconn.Collection("warung").InsertOne(context.Background(), tempatWarung)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Tempat warung berhasil disimpan dengan ID: %s", insertedID.Hex())})
}

func DeleteTempatWarungById(respw http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "ID tidak ditemukan dalam permintaan"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "ID tidak valid"})
		return
	}

	filter := bson.M{"_id": objectID}
	result, err := config.Mongoconn.Collection("warung").DeleteOne(context.Background(), filter)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, itmodel.Response{Response: "Data warung tidak ditemukan"})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: "Data warung berhasil dihapus"})
}

func UpdateTempatWarungById(respw http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "ID tidak ditemukan dalam permintaan"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "ID tidak valid"})
		return
	}

	var updatedData model.Warung
	if err := json.NewDecoder(req.Body).Decode(&updatedData); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "Format JSON tidak valid"})
		return
	}

	if updatedData.Gambar != "" {
		updatedData.Gambar = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + updatedData.Gambar
	}

	updatefields := bson.M{
		"nama_tempat":       updatedData.Nama_Tempat,
		"lokasi":            updatedData.Lokasi,
		"jam_buka":          updatedData.Jam_Buka,
		"metode_pembayaran": updatedData.Metode_Pembayaran,
		"lon":               updatedData.Lon,
		"lat":               updatedData.Lat,
		"gambar":            updatedData.Gambar,
	}

	update := bson.M{
		"$set": updatefields,
	}

	filter := bson.M{"_id": objectID}
	result, err := config.Mongoconn.Collection("warung").UpdateOne(context.Background(), filter, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, itmodel.Response{Response: "Data warung tidak ditemukan"})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: "Data warung berhasil diperbarui"})
}
