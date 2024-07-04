package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/ghupload"
	"github.com/whatsauth/itmodel"
)

func GetGithubFiles(w http.ResponseWriter, r *http.Request) {
	log.Println("GetGithubFiles: Received request")
	var respn itmodel.Response

	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		log.Printf("GetOneDoc error: %v", err)
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	content, err := ghupload.GithubListFiles(gh.GitHubAccessToken, "parkirgratis", "filegambar", "img")
	if err != nil {
		respn.Response = err.Error()
		log.Printf("GithubListFiles error: %v", err)
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	log.Printf("GetGithubFiles: %v", content)

	contentJSON, err := json.Marshal(content)
	if err != nil {
		respn.Response = err.Error()
		log.Printf("json.Marshal error: %v", err)
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = "Files retrieved successfully"
	respn.Response = string(contentJSON)
	helper.WriteJSON(w, http.StatusOK, respn)
}

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

func UpdateGithubFile(w http.ResponseWriter, r *http.Request) {
	var respn itmodel.Response

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()

	// Get the file name from form
	fileName := r.FormValue("fileName")
	if fileName == "" {
		respn.Response = "File name is required"
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	// Get GitHub credentials from the database
	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	// Create a multipart.FileHeader from the uploaded file
	fileHeader := &multipart.FileHeader{
		Filename: handler.Filename,
		Header:   handler.Header,
		Size:     handler.Size,
	}

	// Update the file in GitHub
	content, _, err := ghupload.GithubUpdateFile(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, fileHeader, "parkirgratis", "filegambar", fileName)
	if err != nil {
		respn.Info = "Failed to update GitHub file"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = "File updated successfully"
	respn.Response = *content.Content.Path
	helper.WriteJSON(w, http.StatusOK, respn)
}

func DeleteGithubFile(w http.ResponseWriter, r *http.Request) {
	var respn itmodel.Response
	var deleteRequest struct {
		FileName string `json:"fileName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	_, _, err = ghupload.GithubDeleteFile(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, "parkirgratis", "filegambar", deleteRequest.FileName)
	if err != nil {
		respn.Info = "Failed to delete GitHub file"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = "File deleted successfully"
	respn.Response = deleteRequest.FileName
	helper.WriteJSON(w, http.StatusOK, respn)
}
