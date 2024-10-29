package controller

import (
	"fmt"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/ghupload"
	"github.com/whatsauth/itmodel"
)

func PostUploadGithub(w http.ResponseWriter, r *http.Request) {
	var respn itmodel.Response

	fmt.Println("Starting file upload process")

	_, header, err := r.FormFile("img")
	if err != nil {
		fmt.Println("Error parsing form file:", err)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	folder := helper.GetParam(r)
	var pathFile string
	if folder != "" {
		pathFile = folder + "/" + header.Filename
	} else {
		pathFile = header.Filename
	}

	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		fmt.Println("Error fetching GitHub credentials:", err)
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	content, _, err := ghupload.GithubUpload(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, header, "parkirgratis", "filegambar", pathFile, false)
	if err != nil {
		fmt.Println("Error uploading file to GitHub:", err)
		respn.Info = "gagal upload github"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusEarlyHints, respn)
		return
	}
	if content == nil || content.Content == nil {
		fmt.Println("Error: content or content.Content is nil")
		respn.Response = "Error uploading file"
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = *content.Content.Name
	respn.Response = *content.Content.Path
	helper.WriteJSON(w, http.StatusOK, respn)
	fmt.Println("File upload process completed successfully")
}

func PutUpdateGithub(w http.ResponseWriter, r *http.Request) {
	var respn itmodel.Response

	fmt.Println("Starting file update process")

	
	_, header, err := r.FormFile("img")
	if err != nil {
		fmt.Println("Error parsing form file:", err)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	folder := helper.GetParam(r)
	var pathFile string
	if folder != "" {
		pathFile = folder + "/" + header.Filename
	} else {
		pathFile = header.Filename
	}


	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		fmt.Println("Error fetching GitHub credentials:", err)
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	content, _, err := ghupload.GithubUpload(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, header, "parkirgratis", "filegambar", pathFile, true)
	if err != nil {
		fmt.Println("Error updating file on GitHub:", err)
		respn.Info = "gagal update github"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusEarlyHints, respn)
		return
	}
	if content == nil || content.Content == nil {
		fmt.Println("Error: content or content.Content is nil")
		respn.Response = "Error updating file"
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = *content.Content.Name
	respn.Response = *content.Content.Path
	helper.WriteJSON(w, http.StatusOK, respn)
	fmt.Println("File update process completed successfully")
}

func DeleteFileGithub(w http.ResponseWriter, r *http.Request) {
	var respn itmodel.Response

	fmt.Println("Starting file delete process")

	folder := helper.GetParam(r)
	if folder == "" {
		respn.Response = "File path is required"
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		fmt.Println("Error fetching GitHub credentials:", err)
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	err = ghupload.GithubDelete(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, "parkirgratis", "filegambar", folder)
	if err != nil {
		fmt.Println("Error deleting file from GitHub:", err)
		respn.Info = "gagal menghapus file di github"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Response = "File successfully deleted"
	helper.WriteJSON(w, http.StatusOK, respn)
	fmt.Println("File delete process completed successfully")
}