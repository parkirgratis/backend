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
	
	
		var tempat model.Warung
		err = atdb.FindOneDoc(config.Mongoconn, "warung", bson.M{"_id": objectId}).Decode(&tempat)
		if err != nil {
			helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "Tempat warung not found"})
			return
		}
	
	
		deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "warung", bson.M{"_id": objectId})
		if err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to delete document", "error": err.Error()})
			return
		}
	
		if deleteResult.DeletedCount == 0 {
			helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "Document not found"})
			return
		}
	
		markerId, err := primitive.ObjectIDFromHex("67488d0a8589c79bf4ff6d77")
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
	
		_, err = atdb.UpdateOneArray(config.Mongoconn, "marker_warung", filter, update)
		if err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to update markers", "error": err.Error()})
			return
		}
	
		helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Tempat Warung and markers deleted successfully"})
	}

func PutTempatWarung(respw http.ResponseWriter, req *http.Request) {
	var newWarung model.Warung
	if err := json.NewDecoder(req.Body).Decode(&newWarung); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	if newWarung.ID.IsZero() ||
		newWarung.Nama_Tempat == "" ||
		newWarung.Lokasi == "" ||
		newWarung.Jam_Buka == "" ||
		len(newWarung.Metode_Pembayaran) == 0 ||
		newWarung.Lon == 0 ||
		newWarung.Lat == 0 ||
		newWarung.Gambar == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, "All fields must be filled and valid")
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
    "gambar": newWarung.Gambar,
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
