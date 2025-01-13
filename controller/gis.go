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

func InsertDataRegionFromPetapdia(respw http.ResponseWriter, req *http.Request) {
	var region model.Tempat
	if err := json.NewDecoder(req.Body).Decode(&region); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if region.Province == "" || region.District == "" ||
		region.SubDistrict == "" || region.Village == "" || region.Nama_Tempat == "" ||
		region.Lokasi == "" {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Incomplete region data",
		})
		return
	}

	if region.Gambar != "" {
		region.Gambar = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + region.Gambar
	}

	if region.Lat < -90 || region.Lat > 90 || 
	   region.Lon < -180 || region.Lon > 180 {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Longitude and Latitude must be provided",
		})
		return
	}

	_, err := atdb.InsertOneDoc(config.Mongoconn, "tempat", region)
	if err != nil {
		log.Println("Error saving region to MongoDB:", err)
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
			"error": "Failed to save region to MongoDB",
		})
		return
	}

	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "Success",
		"message": "Region successfully saved",
	})
}

func InsertDataRegionFromPetapdiaWarung(respw http.ResponseWriter, req *http.Request) {
	var region model.Warung
	if err := json.NewDecoder(req.Body).Decode(&region); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if region.Province == "" || region.District == "" ||
		region.SubDistrict == "" || region.Village == "" || region.Nama_Tempat == "" ||
		region.Lokasi == "" || len(region.Metode_Pembayaran) == 0 {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Incomplete region data",
		})
		return
	}

	if region.Gambar != "" {
		region.Gambar = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + region.Gambar
	}

	if region.Lat < -90 || region.Lat > 90 || 
	   region.Lon < -180 || region.Lon > 180 {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
			"error": "Longitude and Latitude must be provided",
		})
		return
	}

	_, err := atdb.InsertOneDoc(config.Mongoconn, "warung", region)
	if err != nil {
		log.Println("Error saving region to MongoDB:", err)
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
			"error": "Failed to save region to MongoDB",
		})
		return
	}

	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "Success",
		"message": "Region successfully saved",
	})
}


func GetRoads(respw http.ResponseWriter, req *http.Request) {
	var longlat model.LongLat
	err := json.NewDecoder(req.Body).Decode(&longlat)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	
	filter := bson.M{
		"geometry": bson.M{
			"$nearSphere": bson.M{ 
				"$geometry": bson.M{ 
					"type":        "Point", 
					"coordinates": []float64{longlat.Longitude, longlat.Latitude},
				},
				"$maxDistance": longlat.MaxDistance,
			},
		},
	}

	var roads []model.Roads
	roads, err = atdb.GetAllDoc[[]model.Roads](config.Mongoconn, "roads", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Tidak ada data ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}
	
	at.WriteJSON(respw, http.StatusOK, roads)
}

func GetRegion(respw http.ResponseWriter, req *http.Request) {
	var longlat model.LongLat
	err := json.NewDecoder(req.Body).Decode(&longlat)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
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

	var region model.Region
	region, err = atdb.GetOneDoc[model.Region](config.Mongoconn, "region", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Region tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, region)
}

func SearchRoadsRegion(respw http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	searchType := queryParams.Get("type") 
	searchText := queryParams.Get("query")

	if searchType == "" || searchText == "" {
		var respn model.Response
		respn.Status = "Error: Parameter 'type' dan 'query' wajib diisi"
		respn.Response = "Parameter 'type' dan 'query' tidak boleh kosong"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var filter bson.M
	var result interface{}

	switch searchType {
	case "roads":
		filter = bson.M{
			"$or": []bson.M{
				{"properties.name": bson.M{"$regex": searchText, "$options": "i"}},
				{"properties.highway": bson.M{"$regex": searchText, "$options": "i"}},
			},
		}
		var roads []model.Roads
		roads, err := atdb.GetAllDoc[[]model.Roads](config.Mongoconn, "roads", filter)
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Data jalan tidak ditemukan"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotFound, respn)
			return
		}
		result = roads

	case "region":
		filter = bson.M{
			"$or": []bson.M{
				{"province": bson.M{"$regex": searchText, "$options": "i"}},
				{"district": bson.M{"$regex": searchText, "$options": "i"}},
				{"sub_district": bson.M{"$regex": searchText, "$options": "i"}},
				{"village": bson.M{"$regex": searchText, "$options": "i"}},
			},
		}
		var regions []model.Region
		regions, err := atdb.GetAllDoc[[]model.Region](config.Mongoconn, "region", filter)
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Data region tidak ditemukan"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotFound, respn)
			return
		}
		result = regions

	default:
		var respn model.Response
		respn.Status = "Error: Tipe pencarian tidak valid"
		respn.Response = "Gunakan 'roads' atau 'region' sebagai tipe pencarian"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, result)
}
