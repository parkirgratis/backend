package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMarker(respw http.ResponseWriter, req *http.Request) {
	mar, err := atdb.GetOneLatestDoc[model.Koordinat](config.Mongoconn, "marker", bson.M{})
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	helper.WriteJSON(respw, http.StatusOK, mar)
}

func PutKoordinat(respw http.ResponseWriter, req *http.Request) {
	var updateRequest struct {
		ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Markers [][]float64        `json:"markers"`
	}

	// Decode request body
	if err := json.NewDecoder(req.Body).Decode(&updateRequest); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Default ID jika tidak diberikan
	id := updateRequest.ID
	if id.IsZero() {
		var err error
		id, err = primitive.ObjectIDFromHex("669510e39590720071a5691d")
		if err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Invalid default ID"})
			return
		}
	}

	
	filter := bson.M{"_id": id}


	doc, err := atdb.GetOneDoc[model.Koordinat](config.Mongoconn, "marker", filter)
	if err != nil {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"error": "Document not found"})
		return
	}

	if len(updateRequest.Markers) < 2 {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Invalid marker data"})
		return
	}


	var index int = -1
	for i, marker := range doc.Markers {
		if len(marker) == 2 && marker[0] == updateRequest.Markers[0][0] && marker[1] == updateRequest.Markers[0][1] {
			index = i
			break
		}
	}

	if index == -1 {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Marker not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("markers.%d", index): updateRequest.Markers[1],
		},
	}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "marker", filter, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Coordinate updated"})
}

func DeleteKoordinat(respw http.ResponseWriter, req *http.Request) {
	var deleteRequest struct {
		ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Markers [][]float64        `json:"markers"`
	}

	// body request
	if err := json.NewDecoder(req.Body).Decode(&deleteRequest); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Validasi format marker
	for _, marker := range deleteRequest.Markers {
		if len(marker) != 2 {
			helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{
				"error": "Invalid marker format, each marker must be an array of [longitude, latitude]",
			})
			return
		}
	}

	// Default ID jika kosong
	id := deleteRequest.ID
	if id.IsZero() {
		var err error
		id, err = primitive.ObjectIDFromHex("669510e39590720071a5691d")
		if err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
				"error": "Invalid default ID",
			})
			return
		}
	}

	// Filter MongoDB
	filter := bson.M{"_id": id}
	update := bson.M{
		"$pull": bson.M{
			"markers": bson.M{
				"$in": deleteRequest.Markers,
			},
		},
	}

	// Update dokumen
	result, err := atdb.UpdateOneDoc(config.Mongoconn, "marker", filter, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{
			"error": "No markers found to delete",
		})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{
		"message": "Coordinates deleted successfully",
	})
}


func PostKoordinat(respw http.ResponseWriter, req *http.Request) {
	var newKoor model.Koordinat
	if err := json.NewDecoder(req.Body).Decode(&newKoor); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	id, err := primitive.ObjectIDFromHex("669510e39590720071a5691d")
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, "Invalid ID format")
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"markers": bson.M{"$each": newKoor.Markers}}}

	if _, err := atdb.UpdateOneDoc(config.Mongoconn, "marker", filter, update); err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(respw, http.StatusOK, "Markers updated")
}