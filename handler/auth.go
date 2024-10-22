package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)
	func GetAdminByUsername(username string) (model.Admin, error) {
		var admin model.Admin

		if config.ErrorMongoconn != nil {
			return admin, fmt.Errorf("failed to connect to database: %w", config.ErrorMongoconn)
		}

		adminCollection := config.Mongoconn.Collection("admin")
		ctx := context.Background()

		err := atdb.FindOne(ctx, adminCollection, bson.M{"username": username}, &admin)
		if err != nil {
			return admin, err
		}

		return admin, nil
	}

	func SaveTokenToMongoWithParams(adminID, token string) error {
		newToken := model.Token{
			AdminID:   adminID,
			Token:     token,
			CreatedAt: time.Now(),
		}

		collection := config.Mongoconn.Collection("tokens")
		ctx := context.Background()

		filter := bson.M{"admin_id": newToken.AdminID}
		update := bson.M{
			"$set": newToken,
		}

		_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		return nil
	}

	func DeleteTokenFromMongo(respw http.ResponseWriter, req *http.Request) error {
		var reqData struct {
			Token string `json:"token"`
		}

		if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
			helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
			return err
		}

		collection := config.Mongoconn.Collection("tokens")
		ctx := context.Background()

		filter := bson.M{"token": reqData.Token}

		// Menghapus token dari database
		_, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to delete token"})
			return err
		}

		helper.WriteJSON(respw, http.StatusOK, map[string]string{"status": "Token deleted successfully"})
		return nil
	}

	func Login(respw http.ResponseWriter, req *http.Request) {
		var loginDetails model.Admin

		if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
			helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
			return
		}

		storedAdmin, err := GetAdminByUsername(loginDetails.Username)
		if err != nil {
			helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Username not found"})
			return
		}

		if loginDetails.Password != storedAdmin.Password {
			helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
			return
		}

		token, err := config.GenerateJWT(storedAdmin.ID.Hex())
		if err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Could not generate token"})
			return
		}

		if err := SaveTokenToMongoWithParams(storedAdmin.ID.Hex(), token); err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Could not save token"})
			return
		}

		helper.WriteJSON(respw, http.StatusOK, map[string]string{
			"status": "Login successful",
			"token":  token,
		})
	}

	func Logout(respw http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Authorization header missing"})
			return
		}
	
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Invalid token format"})
			return
		}
	
		
		collection := config.Mongoconn.Collection("tokens")
		ctx := context.Background()
	
				filter := bson.M{"token": token}
	
		
		_, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to delete token"})
			return
		}
	
		helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Logout successful"})
	}

	func DashboardAdmin(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
	
		adminID := req.Header.Get("admin_id")
		if adminID == "" {
			log.Println("Admin ID tidak ditemukan di header")
			json.NewEncoder(res).Encode(map[string]string{
				"error": "Admin ID tidak ditemukan di header",
			})
			return
		}
	
		resp := map[string]interface{}{
			"status":   http.StatusOK,
			"message":  "Akses dashboard berhasil",
			"admin_id": adminID,
		}
	
		json.NewEncoder(res).Encode(resp)
	}