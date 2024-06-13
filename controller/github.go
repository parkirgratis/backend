package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"

	"github.com/gocroot/helper/ghupload"
	"github.com/gocroot/helper/watoken"
	"github.com/whatsauth/itmodel"
)

func PostUploadGithub(respw http.ResponseWriter, req *http.Request) {
	var respn itmodel.Response
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, helper.GetLoginFromHeader(req))
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(req)
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	// Parse the form file
	_, header, err := req.FormFile("image")
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(req)
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}

	//folder := ctx.Params("folder")
	folder := helper.GetParam(req)
	var pathFile string
	if folder != "" {
		pathFile = folder + "/" + header.Filename
	} else {
		pathFile = header.Filename
	}

	// save to github
	content, _, err := ghupload.GithubUpload(config.GitHubAccessToken, config.GitHubAuthorName, config.GitHubAuthorEmail, header, "parkirgratis.github.io", "release", pathFile, false)
	if err != nil {
		respn.Info = "gagal upload gambar"
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, content)
		return
	}

	helper.WriteJSON(respw, http.StatusOK, respn)

}
