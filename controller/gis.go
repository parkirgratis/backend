package controller

import (
    "encoding/json"
    "net/http"
    "log"

    "github.com/gocroot/config"
    "github.com/gocroot/helper/at"
    "github.com/gocroot/helper/atdb"
    "github.com/gocroot/model"
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

    region := model.Region{
        Province:    "", // Biarkan aja kosong
        District:    "", // Biarkan aja kosong
        SubDistrict: "", // Biarkan aja kosong
        Village:     "", // Biarkan aja kosong
        Border: model.Location{
            Type: "Point", // GeoJSON tipe Point
            Coordinates: [][][]float64{
                {
                    {longlat.Longitude, longlat.Latitude}, // Koordina
                },
            },
        },
    }

    // Simpan data region ke MongoDB
    _, err := atdb.InsertOneDoc(config.Mongoconn, "region", region)
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
        "message": "Region successfully saved to MongoDB",
    })
}