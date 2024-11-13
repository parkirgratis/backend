package controller

import (
	"net/http"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"

)

func GetDataUser(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		docuser.PhoneNumber = payload.Id
		docuser.Name = payload.Alias
		at.WriteJSON(respw, http.StatusNotFound, docuser)
		return
	}
	docuser.Name = payload.Alias
	at.WriteJSON(respw, http.StatusOK, docuser)
}

func PutTokenDataUser(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		docuser.PhoneNumber = payload.Id
		docuser.Name = payload.Alias
		at.WriteJSON(respw, http.StatusNotFound, docuser)
		return
	}
	docuser.Name = payload.Alias
	hcode, qrstat, err := atapi.Get[model.QRStatus](config.WAAPIGetDevice + at.GetLoginFromHeader(req))
	if err != nil {
		at.WriteJSON(respw, http.StatusMisdirectedRequest, docuser)
		return
	}
	if hcode == http.StatusOK && !qrstat.Status {
		docuser.LinkedDevice, err = watoken.EncodeforHours(docuser.PhoneNumber, docuser.Name, config.PrivateKey, 43830)
		if err != nil {
			at.WriteJSON(respw, http.StatusFailedDependency, docuser)
			return
		}
	} else {
		docuser.LinkedDevice = ""
	}
	_, err = atdb.ReplaceOneDoc(config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id}, docuser)
	if err != nil {
		at.WriteJSON(respw, http.StatusExpectationFailed, docuser)
		return
	}
	at.WriteJSON(respw, http.StatusOK, docuser)
}
