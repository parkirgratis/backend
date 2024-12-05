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

	if region.Latitude < -90 || region.Latitude > 90 || 
	   region.Longitude < -180 || region.Longitude > 180 {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Longitude and Latitude must be provided",
		})
		return
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
			"longitude":   region.Longitude,
			"latitude":    region.Latitude,
		},
	})
}
