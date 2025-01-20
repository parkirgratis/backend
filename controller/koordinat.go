package controller

import (
	"encoding/json"
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

func PostKoordinat(respw http.ResponseWriter, req *http.Request) {
	var newKoor model.Koordinat
	if err := json.NewDecoder(req.Body).Decode(&newKoor); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := primitive.ObjectIDFromHex("669510e39590720071a5691d")
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"markers": bson.M{"$each": newKoor.Markers}}}

	_, err = atdb.UpdateOneArray(config.Mongoconn, "marker", filter, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Markers updated"})
}

