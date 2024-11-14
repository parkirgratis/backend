package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/whatsauth/itmodel"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response : err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Tempat warung berhasil disimpan dengan ID: %s", insertedID.Hex())})
}
