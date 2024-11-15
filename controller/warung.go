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
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response : err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Tempat warung berhasil disimpan dengan ID: %s", insertedID.Hex())})
}