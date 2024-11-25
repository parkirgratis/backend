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

    // Membuat struktur Coordinates untuk menyimpan longitude dan latitude
    coordinates := [][][]float64{
        {
            {longlat.Longitude, longlat.Latitude}, // Longitude dan Latitude
        },
    }

    // Membuat objek Location dengan koordinat yang sudah disesuaikan
    location := model.Location{
        Type:        "Point", // Tipe GeoJSON (Point, LineString, Polygon, dll)
        Coordinates: coordinates,
    }

    // Membuat objek Region dan menyimpan lokasi
    region := model.Region{
        Border: location,
    }

    // Menyimpan data ke MongoDB
    _, err := atdb.InsertOneDoc(config.Mongoconn, "region", region)
    if err != nil {
        log.Println("Error saving region to MongoDB:", err)
        at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
            "error": "Failed to save data to MongoDB",
        })
        return
    }

    var response model.Response
    response.Status = "Success"
    response.Response = "Data has been successfully synced and saved to MongoDB"
    at.WriteJSON(respw, http.StatusOK, response)
}
