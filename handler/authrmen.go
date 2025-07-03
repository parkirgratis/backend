package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/controller"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAdminByUsernames(respw http.ResponseWriter, req *http.Request) error {
	var admin model.Admins
	username := req.URL.Query().Get("username")

	if config.ErrorRamenConn != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return fmt.Errorf("failed to connect to database: %w", config.ErrorRamenConn)
	}

	adminCollection := config.RamenConn.Collection("admin")
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

func GetAdminIDFromTokens(respw http.ResponseWriter, req *http.Request) error {
	var admin model.Tokens

	adminID := req.URL.Query().Get("admin_id")
	if adminID == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Admin ID is missing"})
		return fmt.Errorf("admin ID is missing")
	}

	if config.ErrorRamenConn != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
		return fmt.Errorf("failed to connect to database: %w", config.ErrorRamenConn)
	}

	adminCollection := config.RamenConn.Collection("tokens")
	ctx := context.Background()

	// Mencari admin berdasarkan admin_id
	err := atdb.FindOne(ctx, adminCollection, bson.M{"admin_id": adminID}, &admin)
	if err != nil {
		// Jika admin_id tidak ditemukan
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"error": "Admin ID not found"})
		return fmt.Errorf("admin ID not found: %w", err)
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status": "Admin ID found",
		"admin":  admin,
	})
	return nil
}

func SaveTokenToMongoWithParamss(respw http.ResponseWriter, req *http.Request) error {
	var reqData struct {
		AdminID string `json:"admin_id"`
		Token   string `json:"token"`
	}

	// Parsing body JSON untuk mendapatkan admin_id dan token
	if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return err
	}

	// Memastikan admin_id dan token tidak kosong
	if reqData.AdminID == "" || reqData.Token == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Admin ID or Token is missing"})
		return fmt.Errorf("admin ID or token is missing")
	}

	// Membuat token baru untuk penyimpanan
	newToken := model.Tokens{
		AdminID:   reqData.AdminID,
		Token:     reqData.Token,
		CreatedAt: time.Now(),
	}

	collection := config.RamenConn.Collection("tokens")
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

func DeleteTokenFromMongos(respw http.ResponseWriter, req *http.Request) error {
	var reqData struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(req.Body).Decode(&reqData); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		return err
	}

	collection := config.RamenConn.Collection("tokens")
	ctx := context.Background()

	filter := bson.M{"token": reqData.Token}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Failed to delete token"})
		return err
	}

	// Kirim respons berhasil
	helper.WriteJSON(respw, http.StatusOK, map[string]string{"status": "Token deleted successfully"})
	return nil
}

func Logins(respw http.ResponseWriter, req *http.Request) {
	var loginDetails model.Admins
	if err := json.NewDecoder(req.Body).Decode(&loginDetails); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	// 1. Validasi input kosong
	if strings.TrimSpace(loginDetails.Username) == "" || strings.TrimSpace(loginDetails.Password) == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Username and password are required"})
		return
	}

	// 2. Validasi panjang username dan password
	if len(loginDetails.Username) < 3 || len(loginDetails.Username) > 30 {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Username must be between 3 and 30 characters"})
		return
	}

	if len(loginDetails.Password) < 6 || len(loginDetails.Password) > 50 {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Password must be between 6 and 50 characters"})
		return
	}

	// 3. Validasi username hanya mengandung huruf dan angka
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !usernameRegex.MatchString(loginDetails.Username) {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid username format"})
		return
	}

	// 4. Cek apakah akun sudah terlalu banyak gagal login (Rate Limiting)
	if failedAttemptsExceeded(loginDetails.Username) {
		helper.WriteJSON(respw, http.StatusTooManyRequests, map[string]string{"message": "Too many failed login attempts, please try again later"})
		return
	}

	var storedAdmin model.Admins
	if err := atdb.FindOne(context.Background(), config.RamenConn.Collection("admin"), bson.M{"username": loginDetails.Username}, &storedAdmin); err != nil {
		incrementFailedAttempts(loginDetails.Username) // Tambah ke gagal login
		helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
		return
	}

	if !config.CheckPasswordHash(loginDetails.Password, storedAdmin.Password) {
		incrementFailedAttempts(loginDetails.Username) // Tambah ke gagal login
		helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
		return
	}

	// Reset counter gagal login jika berhasil
	resetFailedAttempts(loginDetails.Username)

	// Menambahkan role pada JWT
	token, err := config.GenerateJWTs(storedAdmin.ID.Hex(), storedAdmin.Role)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Could not generate token"})
		return
	}

	collection := config.RamenConn.Collection("tokens")
	ctx := context.Background()
	newToken := model.Tokens{
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

var loginAttempts = make(map[string]int)
var attemptTimestamps = make(map[string]time.Time)

const maxFailedAttempts = 5
const lockoutDuration = 5 * time.Minute

// Cek apakah user telah melewati batas percobaan login
func failedAttemptsExceeded(username string) bool {
	attempts, exists := loginAttempts[username]
	if !exists {
		return false
	}

	// Cek apakah user masih dalam masa blokir
	if attempts >= maxFailedAttempts {
		lastAttempt, _ := attemptTimestamps[username]
		if time.Since(lastAttempt) < lockoutDuration {
			return true
		}
		// Reset jika sudah lewat waktu blokir
		delete(loginAttempts, username)
		delete(attemptTimestamps, username)
	}

	return false
}

// Tambah jumlah percobaan gagal login
func incrementFailedAttempts(username string) {
	loginAttempts[username]++
	attemptTimestamps[username] = time.Now()
}

// Reset percobaan gagal jika login berhasil
func resetFailedAttempts(username string) {
	delete(loginAttempts, username)
	delete(attemptTimestamps, username)
}

func Logouts(respw http.ResponseWriter, req *http.Request) {
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

	collection := config.RamenConn.Collection("tokens")
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

func DashboardAdmins(res http.ResponseWriter, req *http.Request) {
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

func RegisterAdmins(respw http.ResponseWriter, req *http.Request) {
	var adminDetails model.Admins

	// Decode the request body into adminDetails struct
	if err := json.NewDecoder(req.Body).Decode(&adminDetails); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	// Validate input fields
	if adminDetails.Username == "" || adminDetails.Password == "" || adminDetails.Role == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Username, password, and role are required"})
		return
	}

	// Validate username format (alphanumeric and no spaces)
	if !isValidUsername(adminDetails.Username) {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Username must be alphanumeric and cannot contain spaces"})
		return
	}

	// Validate password strength (minimum 6 characters, must contain at least one number and one letter)
	if !isValidPassword(adminDetails.Password) {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Password must be at least 6 characters long and contain at least one number and one letter"})
		return
	}

	// Validate role (role must be 'admin' or 'kasir')
	if !isValidRole(adminDetails.Role) {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Role must be 'admin' or 'kasir'"})
		return
	}

	// Check if username already exists in the database
	var existingAdmin model.Admins
	err := atdb.FindOne(context.Background(), config.RamenConn.Collection("admin"), bson.M{"username": adminDetails.Username}, &existingAdmin)
	if err == nil {
		helper.WriteJSON(respw, http.StatusConflict, map[string]string{"message": "Username already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := config.HashPassword(adminDetails.Password)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
		return
	}

	// Create new admin object
	newAdmin := model.Admins{
		Username: adminDetails.Username,
		Password: hashedPassword,
		Role:     adminDetails.Role,
	}

	// Insert the new admin into the database
	collection := config.RamenConn.Collection("admin")
	ctx := context.Background()

	_, err = collection.InsertOne(ctx, newAdmin)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to register admin"})
		return
	}

	// Return success response
	helper.WriteJSON(respw, http.StatusCreated, map[string]string{
		"status":   "Admin registered successfully",
		"username": newAdmin.Username,
	})
}

// Helper function to validate username
func isValidUsername(username string) bool {
	// Username should be alphanumeric and should not contain spaces
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return re.MatchString(username)
}

// Helper function to validate password strength
func isValidPassword(password string) bool {
	// Password must be at least 6 characters long and contain at least one letter and one number
	if len(password) < 6 {
		return false
	}
	hasLetter := false
	hasNumber := false
	for _, char := range password {
		if strings.ContainsAny(string(char), "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			hasLetter = true
		}
		if strings.ContainsAny(string(char), "0123456789") {
			hasNumber = true
		}
	}
	return hasLetter && hasNumber
}

// Helper function to validate role
func isValidRole(role string) bool {
	// Validate that the role is either "admin" or "kasir"
	return role == "admin" || role == "kasir"
}

func UpdateForgottenPassword(respw http.ResponseWriter, req *http.Request) {
	var updateRequest struct {
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}

	// Decode the request body into updateRequest struct
	if err := json.NewDecoder(req.Body).Decode(&updateRequest); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	// Validate input fields
	if updateRequest.Username == "" || updateRequest.NewPassword == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Username and new password are required"})
		return
	}

	// Validate new password strength (minimum 6 characters, must contain at least one number and one letter)
	if !isValidPassword(updateRequest.NewPassword) {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "New password must be at least 6 characters long and contain at least one number and one letter"})
		return
	}

	// Check if username exists in the database
	var existingAdmin model.Admins
	collection := config.RamenConn.Collection("admin")
	ctx := context.Background()

	err := collection.FindOne(ctx, bson.M{"username": updateRequest.Username}).Decode(&existingAdmin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "Username not found"})
		} else {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to find admin"})
		}
		return
	}

	// Hash the new password
	hashedPassword, err := config.HashPassword(updateRequest.NewPassword)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to hash new password"})
		return
	}

	// Update the password in the database
	update := bson.M{"$set": bson.M{"password": hashedPassword}}
	_, err = collection.UpdateOne(ctx, bson.M{"username": updateRequest.Username}, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to update password"})
		return
	}

	// Return success response
	helper.WriteJSON(respw, http.StatusOK, map[string]string{
		"status":   "Password updated successfully",
		"username": updateRequest.Username,
	})
}

func GetAllAdmins(respw http.ResponseWriter, req *http.Request) {
	// Mendapatkan koneksi ke collection "admin"
	collection := config.RamenConn.Collection("admin")
	ctx := context.Background()

	// Mencari semua data admin dalam collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to fetch admin data"})
		return
	}
	defer cursor.Close(ctx)

	// Membuat slice untuk menyimpan hasil query
	var admins []model.Admins

	// Iterasi melalui cursor dan decode setiap dokumen ke dalam slice admins
	for cursor.Next(ctx) {
		var admin model.Admins
		if err := cursor.Decode(&admin); err != nil {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to decode admin data"})
			return
		}
		admins = append(admins, admin)
	}

	// Periksa apakah ada error selama iterasi
	if err := cursor.Err(); err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Error during cursor iteration"})
		return
	}

	// Jika tidak ada data admin yang ditemukan
	if len(admins) == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "No admin data found"})
		return
	}

	// Mengembalikan data admin dalam bentuk JSON
	helper.WriteJSON(respw, http.StatusOK, admins)
}
