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

	if tempatWarung.Foto_pratinjau != "" {
		tempatWarung.Foto_pratinjau = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + tempatWarung.Foto_pratinjau
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
	id := req.URL.Query().Get("_id")
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

func PutTempatWarung(respw http.ResponseWriter, req *http.Request) {
	var newWarung model.Warung
	if err := json.NewDecoder(req.Body).Decode(&newWarung); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Decoded document:", newWarung)
	id, err := primitive.ObjectIDFromHex(newWarung.ID.Hex())
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, "Invalid ID format")
		return
	}

	filter := bson.M{"_id": id}
	updatefields := bson.M{
    "nama_tempat": newWarung.Nama_Tempat,
    "lokasi": newWarung.Lokasi,
	"jam_buka": newWarung.Jam_Buka,
	"metode_pembayaran": newWarung.Metode_Pembayaran,
    "lon": newWarung.Lon,
    "lat": newWarung.Lat,
    "foto_pratinjau": newWarung.Foto_pratinjau,
}

	fmt.Println("Filter:", filter)
	fmt.Println("Update:", updatefields)

	result, err := atdb.UpdateOneDoc(config.Mongoconn, "warung", filter, updatefields)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, "Document not found or not modified")
		return
	}


	helper.WriteJSON(respw, http.StatusOK, newWarung)
}
