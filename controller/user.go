package controller

import (
	"net/http"

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