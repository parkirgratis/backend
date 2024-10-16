package handler

import(
	"context"
	"encoding/json"
	"fmt"

	"time"
	"net/http"

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
		http.Error(respw, "Invalid request body", http.StatusBadRequest)
		return
	}

	storedAdmin, err := GetAdminByUsername(loginDetails.Username)
	if err != nil {
		helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Username not found"})
		return
	}

	if loginDetails.Password != storedAdmin.Password {
		http.Error(respw, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := config.GenerateJWT(storedAdmin.ID.Hex())
	if err != nil {
		http.Error(respw, "Could not generate token", http.StatusInternalServerError)
		return
	}

	if err := SaveTokenToMongoWithParams(respw, req); err != nil {
		http.Error(respw, "Could not save token", http.StatusInternalServerError)
		return
	}

	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	json.NewEncoder(respw).Encode(map[string]string{
		"status": "Login successful",
		"token":  token,
	})
}

func DashboardAdmin(respw http.ResponseWriter, req *http.Request) {
	adminID := req.Context().Value("admin_id")
	if adminID == nil {
		http.Error(respw, "Admin ID not found in context", http.StatusInternalServerError)
		return
	}

	adminIDStr := fmt.Sprintf("%v", adminID)

	respw.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"status":   http.StatusOK,
		"message":  "Dashboard access successful",
		"admin_id": adminIDStr,
	}
	json.NewEncoder(respw).Encode(resp)
}