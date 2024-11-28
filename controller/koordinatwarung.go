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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func GetMarkerWarung(respw http.ResponseWriter, req *http.Request) {
	mar, err := atdb.GetOneLatestDoc[model.KoordinatWarung](config.Mongoconn, "marker_warung", bson.M{})
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	helper.WriteJSON(respw, http.StatusOK, mar)
}

func PutKoordinatWarung(respw http.ResponseWriter, req *http.Request) {
	var updateRequest struct {
		ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Markers [][]float64        `json:"markers"`
	}

	if err := json.NewDecoder(req.Body).Decode(&updateRequest); err != nil {
		http.Error(respw, err.Error(), http.StatusBadRequest)
		return
	}

	id := updateRequest.ID
	if id.IsZero() {
		defaultID, err := primitive.ObjectIDFromHex("669510e39590720071a5691d")
		if err != nil {
			http.Error(respw, "Invalid default ID", http.StatusInternalServerError)
			return
		}
		id = defaultID
	}

	filter := bson.M{"_id": id}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://irgifauzi:%40Sasuke123@webservice.rq9zk4m.mongodb.net/"))
	if err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("parkir_db").Collection("marker_warung")

	var document struct {
		Markers [][]float64 `bson:"markers"`
	}
	if err := collection.FindOne(context.TODO(), filter).Decode(&document); err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}

	var index int = -1
	for i, marker := range document.Markers {
		if len(marker) == 2 && marker[0] == updateRequest.Markers[0][0] && marker[1] == updateRequest.Markers[0][1] {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(respw, "Marker not found", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("markers.%d", index): updateRequest.Markers[1],
		},
	}

	if _, err := collection.UpdateOne(context.TODO(), filter, update); err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}

	respw.WriteHeader(http.StatusOK)
	respw.Write([]byte("Coordinate updated"))
}

func DeleteKoordinatWarung(respw http.ResponseWriter, req *http.Request) {
	var deleteRequest struct {
		ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Markers [][]float64        `json:"markers"`
	}

	if err := json.NewDecoder(req.Body).Decode(&deleteRequest); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	id := deleteRequest.ID

	filter := bson.M{"_id": id}
	update := bson.M{
		"$pull": bson.M{
			"markers": bson.M{
				"$in": deleteRequest.Markers,
			},
		},
	}

	result, err := atdb.UpdateOneDoc(config.Mongoconn, "marker_warung", filter, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, "No markers found to delete")
		return
	}

	helper.WriteJSON(respw, http.StatusOK, "Coordinates deleted")
}
