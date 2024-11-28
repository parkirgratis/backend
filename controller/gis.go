package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
)

func SyncDataWithPetapedia(respw http.ResponseWriter, req *http.Request) {
    var longlat model.LongLat
    if err := json.NewDecoder(req.Body).Decode(&longlat); err != nil || longlat.Latitude == 0 || longlat.Longitude == 0 {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
            "error": "Invalid latitude or longitude",
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

    region, err := atdb.GetOneDoc[model.Region](config.Mongoconn, "region", filter)
    if err != nil {
        at.WriteJSON(respw, http.StatusNotFound, map[string]string{
            "error": "Region not found in Petapedia",
        })
        return
    }

    parkirRegion := bson.M{
        "longitude":    longlat.Longitude,
        "latitude":     longlat.Latitude,
        "province":     region.Province,
        "district":     region.District,
        "sub_district":  region.SubDistrict,
        "village":      region.Village,
        "border":       region.Border, 
    }

    _, err = atdb.InsertOneDoc(config.Mongoconn, "parkir_regions", parkirRegion)
    if err != nil {
        log.Println("Error saving region to MongoDB:", err)
        at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
            "error": "Failed to save region to MongoDB",
        })
        return
    }

    at.WriteJSON(respw, http.StatusOK, map[string]string{
        "status":  "Success",
        "message": "Region successfully synced from Petapedia and saved to MongoDB",
    })
}
