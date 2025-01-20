package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/gocroot/controller"
	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func GetAdminByUsername(respw http.ResponseWriter, req *http.Request) error {
	var admin model.Admin
	username := req.URL.Query().Get("username") 

	if config.ErrorMongoconn != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return fmt.Errorf("failed to connect to database: %w", config.ErrorMongoconn)
	}

	adminCollection := config.Mongoconn.Collection("admin")
	ctx := context.Background()

	err := atdb.FindOne(ctx, adminCollection, bson.M{"username": username}, &admin)
	if err != nil {
		
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"error": "User not found"})
		return err
	}

	
	helper.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status": "Admin found",
		"admin":  admin,
	})
	return nil
}


func GetAdminIDFromToken(respw http.ResponseWriter, req *http.Request) error {
	var admin model.Token

	adminID := req.URL.Query().Get("admin_id")
	if adminID == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Admin ID is missing"})
		return fmt.Errorf("admin ID is missing")
	}

	if config.ErrorMongoconn != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return fmt.Errorf("failed to connect to database: %w", config.ErrorMongoconn)
	}

	adminCollection := config.Mongoconn.Collection("tokens") 
	ctx := context.Background()

	// Mencari admin berdasarkan admin_id
	err := atdb.FindOne(ctx, adminCollection, bson.M{"admin_id": adminID}, &admin)
	if err != nil {
		// Jika admin_id tidak ditemukan
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"error": "Admin ID not found"})
		return fmt.Errorf("admin ID not found: %w", err)
	}

	// Jika ditemukan, kirim respons berhasil
	helper.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status": "Admin ID found",
		"admin":  admin,
	})
	return nil
}

// SaveTokenToMongoWithParams menyimpan token ke MongoDB dengan parameter yang diberikan
func SaveTokenToMongoWithParams(respw http.ResponseWriter, req *http.Request) error {
	var reqData struct {
		AdminID string `json:"admin_id"`
		Token   string `json:"token"`
	}

	if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return err
	}

	if reqData.AdminID == "" || reqData.Token == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Admin ID or Token is missing"})
		return fmt.Errorf("admin ID or token is missing")
	}

	newToken := model.Token{
		AdminID:   reqData.AdminID,
		Token:     reqData.Token,
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
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to save token"})
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

    var storedAdmin model.Admin
    if err := atdb.FindOne(context.Background(), config.Mongoconn.Collection("admin"), bson.M{"username": loginDetails.Username}, &storedAdmin); err != nil {
        helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
        return
    }

    if !config.CheckPasswordHash(loginDetails.Password, storedAdmin.Password) {
        helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
        return
    }

    token, err := config.GenerateJWT(storedAdmin.ID.Hex())
    if err != nil {
        helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Could not generate token"})
        return
    }

    collection := config.Mongoconn.Collection("tokens")
    ctx := context.Background()
    newToken := model.Token{
        AdminID:   storedAdmin.ID.Hex(),
        Token:     token,
        CreatedAt: time.Now(),
    }

    _, err = collection.InsertOne(ctx, newToken)
    if err != nil {
        helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Could not save token"})
        return
    }

    if err := controller.LogActivity(storedAdmin.Username); err != nil {
        helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to log login activity"})
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

	// Kirim respons berhasil
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

func RegisterAdmin(respw http.ResponseWriter, req *http.Request) {
	var adminDetails model.Admin
	
	if err := json.NewDecoder(req.Body).Decode(&adminDetails); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	
	if adminDetails.Username == "" || adminDetails.Password == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Username and password are required"})
		return
	}

	
	var existingAdmin model.Admin
	err := atdb.FindOne(context.Background(), config.Mongoconn.Collection("admin"), bson.M{"username": adminDetails.Username}, &existingAdmin)
	if err == nil {
		helper.WriteJSON(respw, http.StatusConflict, map[string]string{"message": "Username already exists"})
		return
	}


	hashedPassword, err := config.HashPassword(adminDetails.Password)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
		return
	}


	newAdmin := model.Admin{
		Username: adminDetails.Username,
		Password: hashedPassword,
	}

	
	collection := config.Mongoconn.Collection("admin")
	ctx := context.Background()

	
	_, err = collection.InsertOne(ctx, newAdmin)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to register admin"})
		return
	}

	
	helper.WriteJSON(respw, http.StatusCreated, map[string]string{
		"status":   "Admin registered successfully",
		"username": newAdmin.Username,
	})
}
