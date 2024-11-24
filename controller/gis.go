package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gocroot/helper/at"
	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"github.com/gocroot/helper/atdb"
)

func SyncDataWithPetapedia(respw http.ResponseWriter, req *http.Request) {
	var coordinates struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	}

	err := json.NewDecoder(req.Body).Decode(&coordinates)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return
	}

	petapediaAPI := "https://asia-southeast2-awangga.cloudfunctions.net/petabackend/data/gis/lokasi"

	requestBody, err := json.Marshal(coordinates)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to encode request body"})
		return
	}

	reqPetapedia, err := http.NewRequest("POST", petapediaAPI, bytes.NewBuffer(requestBody))
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to create request for Petapedia"})
		return
	}
	reqPetapedia.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respPetapedia, err := client.Do(reqPetapedia)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadGateway, map[string]string{"error": "Failed to send request to Petapedia"})
		return
	}
	defer respPetapedia.Body.Close()

	if respPetapedia.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(respPetapedia.Body)
		at.WriteJSON(respw, http.StatusBadGateway, map[string]string{"error": "Petapedia API returned an error", "details": string(body)})
		return
	}

	var region model.Region
	err = json.NewDecoder(respPetapedia.Body).Decode(&region)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to decode Petapedia response"})
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "region", region)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to save data to MongoDB"})
		return
	}

	var response model.Response
	response.Status = "Success"
	response.Response = "Data has been successfully synced with Petapedia and saved to MongoDB"
	at.WriteJSON(respw, http.StatusOK, response)
}
