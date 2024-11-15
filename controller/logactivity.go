package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LogActivity(respw http.ResponseWriter, req *http.Request) error {

	adminID := req.Header.Get("admin_id")
	if adminID == "" {
		http.Error(respw, "Admin ID not found", http.StatusUnauthorized)
		return fmt.Errorf("admin ID missing")
	}

	userID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		http.Error(respw, "Invalid Admin ID", http.StatusBadRequest)
		return fmt.Errorf("invalid admin ID: %v", err)
	}

	logactivity := struct {
		UserID    primitive.ObjectID `bson:"admin_id,omitempty" json:"admin_id,omitempty"`
		Action    string             `bson:"action,omitempty" json:"action,omitempty"`
		Timestamp time.Time          `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	}{
		UserID:    userID,
		Action:    "Your Action",
		Timestamp: time.Now(),
	}

	collection := config.Mongoconn.Collection("activity_logs")
	_, err = collection.InsertOne(context.Background(), logactivity)
	if err != nil {
		http.Error(respw, "Failed to log activity", http.StatusInternalServerError)
		return err
	}

	return nil
}