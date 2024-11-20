package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveLocation(respw http.ResponseWriter, req *http.Request) {
	var location model.Location
	if err := json.NewDecoder(req.Body).Decode(&location); err != nil {
		http.Error(respw, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	location.CreatedAt = time.Now()

	result, err := config.Mongoconn.Collection("locations").InsertOne(context.Background(), location)
	if err != nil {
		http.Error(respw, "Failed to save location: "+err.Error(), http.StatusInternalServerError)
		return
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		http.Error(respw, "Failed to cast InsertedID to ObjectID", http.StatusInternalServerError)
		return
	}

	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	json.NewEncoder(respw).Encode(map[string]interface{}{
		"status":      "success",
		"message":     "Location saved successfully",
		"data":        location,
		"inserted_id": insertedID.Hex(),
	})
}
