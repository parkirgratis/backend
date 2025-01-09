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
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetTempatByNamaTempat(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	lokasi := req.URL.Query().Get("nama_tempat")

	filter := bson.M{"nama_tempat": bson.M{"$regex": lokasi, "$options": "i"}}
	opts := options.Find().SetLimit(10)

	tempat, err := atdb.GetFilteredDocs[[]model.Tempat](config.Mongoconn, "tempat", filter, opts)
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}

	helper.WriteJSON(respw, http.StatusOK, tempat)
}

func PostTempatParkir(respw http.ResponseWriter, req *http.Request) {

	var tempatParkir model.Tempat
	if err := json.NewDecoder(req.Body).Decode(&tempatParkir); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	if tempatParkir.Gambar != "" {
		tempatParkir.Gambar = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + tempatParkir.Gambar
	}

	result, err := config.Mongoconn.Collection("tempat").InsertOne(context.Background(), tempatParkir)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Tempat parkir berhasil disimpan dengan ID: %s", insertedID.Hex())})
}

func PutTempatParkir(respw http.ResponseWriter, req *http.Request) {
	var newTempat model.Tempat
	if err := json.NewDecoder(req.Body).Decode(&newTempat); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	if newTempat.ID.IsZero() {
		helper.WriteJSON(respw, http.StatusBadRequest, "ID is required")
		return
	}

	if newTempat.Gambar != "" {
		newTempat.Gambar = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + newTempat.Gambar
	}

	filter := bson.M{"_id": newTempat.ID}
	updatefields := bson.M{
		"nama_tempat": newTempat.Nama_Tempat,
		"lokasi":      newTempat.Lokasi,
		"fasilitas":   newTempat.Fasilitas,
		"lon":         newTempat.Lon,
		"lat":         newTempat.Lat,
		"gambar":      newTempat.Gambar,
	}

	result, err := atdb.UpdateOneDoc(config.Mongoconn, "tempat", filter, updatefields)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, "Document not found or not modified")
		return
	}

	helper.WriteJSON(respw, http.StatusOK, newTempat)
}

//fungsi delete parkir dan marker
func DeleteTempatParkir(respw http.ResponseWriter, req *http.Request) {
	var requestBody struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid JSON data"})
		return
	}

	
	objectId, err := primitive.ObjectIDFromHex(requestBody.ID)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid ID format"})
		return
	}


	var tempat model.Tempat
	err = atdb.FindOneDoc(config.Mongoconn, "tempat", bson.M{"_id": objectId}).Decode(&tempat)
	if err != nil {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "Tempat not found"})
		return
	}


	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "tempat", bson.M{"_id": objectId})
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to delete document", "error": err.Error()})
		return
	}

	if deleteResult.DeletedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "Document not found"})
		return
	}

	markerId, err := primitive.ObjectIDFromHex("669510e39590720071a5691d")
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Invalid ObjectID format", "error": err.Error()})
		return
	}

	filter := bson.M{"_id": markerId}
	update := bson.M{
		"$pull": bson.M{
			"markers": []float64{tempat.Lon, tempat.Lat},
		},
	}

	_, err = atdb.UpdateOneArray(config.Mongoconn, "marker", filter, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to update markers", "error": err.Error()})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Tempat and markers deleted successfully"})
}
