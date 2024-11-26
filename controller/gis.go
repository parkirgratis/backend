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
    // Validasi dan decode body menjadi satu struct
    var locationData model.LocationData
    if err := json.NewDecoder(req.Body).Decode(&locationData); err != nil || locationData.Latitude == 0 || locationData.Longitude == 0 {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
            "error": "Invalid latitude or longitude",
        })
        return
    }

    // Membuat region dan menambahkan data lokasi
    region := model.Region{
        Province:    locationData.Region.Province,
        District:    locationData.Region.District, 
        SubDistrict: locationData.Region.SubDistrict,
        Village:     locationData.Region.Village,
        Border: model.Location{
            Type: "Point",
            Coordinates: [][][]float64{
                {
                    {locationData.Longitude, locationData.Latitude},
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