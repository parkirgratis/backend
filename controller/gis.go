package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"log"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"github.com/gocroot/helper/watoken"
)

func SyncDataWithPetapedia(respw http.ResponseWriter, req *http.Request) {
	//validate watoken
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		at.WriteJSON(respw, http.StatusForbidden, map[string]string{
			"error": "Invalid token",
		})
		return
	}

	var longlat model.LongLat
	err = json.NewDecoder(req.Body).Decode(&longlat)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
		return
	}

	if longlat.Latitude == 0 || longlat.Longitude == 0 {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Invalid latitude or longitude",
		})
		return
	}

	petapediaAPI := "https://asia-southeast2-awangga.cloudfunctions.net/petabackend/data/gis/lokasi"

	requestBody, err := json.Marshal(longlat)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
			"error": "Failed to encode request body",
		})
		return
	}

	reqPetapedia, err := http.NewRequest("POST", petapediaAPI, bytes.NewBuffer(requestBody))
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create request for Petapedia",
		})
		return
	}
	reqPetapedia.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respPetapedia, err := client.Do(reqPetapedia)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadGateway, map[string]string{
			"error": "Failed to send request to Petapedia",
		})
		return
	}
	defer respPetapedia.Body.Close()

	if respPetapedia.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respPetapedia.Body)
		log.Println("Petapedia API error response:", string(body))
		at.WriteJSON(respw, http.StatusBadGateway, map[string]string{
			"error":   "Petapedia API returned an error",
			"details": string(body),
		})
		return
	}

	var region model.Region
	err = json.NewDecoder(respPetapedia.Body).Decode(&region)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
			"error": "Failed to decode Petapedia response",
		})
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "region", region)
	if err != nil {
		log.Println("Error saving region to MongoDB:", err)
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
			"error": "Failed to save data to MongoDB",
		})
		return
	}

	var response model.Response
	response.Status = "Success"
	response.Response = "Data has been successfully synced with Petapedia and saved to MongoDB"
	at.WriteJSON(respw, http.StatusOK, response)
}
