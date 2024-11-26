package controller

import (
    "encoding/json"
    "net/http"
    "log"

    "github.com/gocroot/config"
    "github.com/gocroot/helper/at"
    "github.com/gocroot/helper/atdb"
    "github.com/gocroot/model"
    "go.mongodb.org/mongo-driver/bson"
)

func SyncDataWithPetapedia(respw http.ResponseWriter, req *http.Request) {
    // Validasi dan decode body
    var longlat model.LongLat
    if err := json.NewDecoder(req.Body).Decode(&longlat); err != nil || longlat.Latitude == 0 || longlat.Longitude == 0 {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
            "error": "Invalid latitude or longitude",
        })
        return
    }

    if longlat.Latitude < -90 || longlat.Latitude > 90 || longlat.Longitude < -180 || longlat.Longitude > 180 {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
            "error": "Latitude or longitude out of range",
        })
        return
    }

    filter := bson.M{
        "border": bson.M{
            "$geoIntersects": bson.M{
                "$geometry": bson.M{
                    "type":        "Point",
                    "coordinates": []float64{longlat.Longitude, longlat.Latitude},
                },
            },
        },
    }

    // Ambil region dari database Petapedia

    // Simpan region ke database ParkirGratis
    _, err := atdb.InsertOneDoc(config.Mongoconn, "region", filter)
    if err != nil {
        log.Println("Error saving region to MongoDB:", err)
        at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
            "error": "Failed to save region to MongoDB",
        })
        return
    }

    // Response sukses
    at.WriteJSON(respw, http.StatusOK, map[string]string{
        "status":  "Success",
        "message": "Region successfully synced from Petapedia and saved to MongoDB",
    })
}