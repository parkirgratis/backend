package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// hahay
func LogActivityss(username string) error {
	collection := config.Mongoconn.Collection("activity_logs")
	ctx := context.Background()

	log := model.LoginLogs{
		ID:        primitive.NewObjectID().Hex(),
		Username:  username,
		Activity:  "Login",
		Timestamp: time.Now(),
	}

	_, err := collection.InsertOne(ctx, log)
	return err
}

func GetActivity(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response

	activities, err := atdb.GetAllDoc[[]model.LoginLogs](config.Mongoconn, "activity_logs", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}

	helper.WriteJSON(respw, http.StatusOK, activities)
}
