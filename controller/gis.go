package controller

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"github.com/gocroot/helper/atdb"
)

func SyncDataFromPetapediaAPI(respw http.ResponseWriter, req *http.Request) {
    petapediadAPI := "https://asia-southeast2-awangga.cloudfunctions.net/petabackend/data/gis/lokasi"
    resp, err := http.Get(petapediadAPI)
    if err != nil {
        var respn model.Response
        respn.Status = "Error : Tidak dapat mengakses API teman"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusInternalServerError, respn)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        var respn model.Response
        respn.Status = "Error : Response API teman gagal"
        respn.Response = fmt.Sprintf("Status code: %d", resp.StatusCode)
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    var locations []model.Region 
    err = json.NewDecoder(resp.Body).Decode(&locations)
    if err != nil {
        var respn model.Response
        respn.Status = "Error : Gagal decode data dari API teman"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusInternalServerError, respn)
        return
    }

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

    var respn model.Response
    respn.Status = "Sukses"
    respn.Response = "Data berhasil disinkronkan dari API teman"
    at.WriteJSON(respw, http.StatusOK, respn)
}

