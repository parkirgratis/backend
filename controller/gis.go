package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
)

func InsertDataRegionFromPetapdia(respw http.ResponseWriter, req *http.Request) {
	var region model.Region

	if err := json.NewDecoder(req.Body).Decode(&region); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if region.Province == "" || region.District == "" ||
		region.SubDistrict == "" || region.Village == "" {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Incomplete region data",
		})
		return
	}

	if region.Longitude == 0 || region.Latitude == 0 {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Longitude and Latitude must be provided",
		})
		return
	}

	if len(region.Border.Coordinates) == 0 {
		region.Border = model.Location{
			Type:        "Point",
			Coordinates: [][][]float64{},
		}
	}

	longLat := model.LongLat{
		Longitude: region.Longitude,
		Latitude:  region.Latitude,
	}

	_, err := atdb.InsertOneDoc(config.Mongoconn, "region", region)
	if err != nil {
		log.Println("Error saving region to MongoDB:", err)
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
			"error": "Failed to save region to MongoDB",
		})
		return
	}

	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "Success",
		"message": "Region successfully saved to MongoDB",
		"data": map[string]interface{}{
			"province":    region.Province,
			"district":    region.District,
			"subDistrict": region.SubDistrict,
			"village":     region.Village,
			"longitude":   longLat.Longitude,
			"latitude":    longLat.Latitude,
		},
	})
}
