package controller

import (
	"encoding/json"
	"net/http"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"github.com/gocroot/helper/atdb"
)

func SyncDataPetapediaBackend(respw http.ResponseWriter, req *http.Request) {
    var locations []model.Region
	err := json.NewDecoder(req.Body).Decode(&locations)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Data yang diterima tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Menyimpan data ke database
	for _, location := range locations {
		_, err := atdb.InsertOneDoc(config.Mongoconn, "region", location)
		if err != nil {
			var respn model.Response
			respn.Status = "Error : Gagal menyimpan data ke database"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusInternalServerError, respn)
			return
		}
	}

	// Mengirimkan respon sukses
	var respn model.Response
	respn.Status = "Sukses"
	respn.Response = "Data berhasil disimpan ke database"
	at.WriteJSON(respw, http.StatusOK, respn)
}

