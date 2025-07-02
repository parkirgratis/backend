package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"

	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// hahay
func GetMenu_ramen(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	resto, err := atdb.GetAllDoc[[]model.Menus](config.Mongoconn, "menu_ramen", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, resto)
}
func GetMenu_ramenflutter(respw http.ResponseWriter, req *http.Request) {
	var resp struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	resto, err := atdb.GetAllDoc[[]model.Menu](config.Mongoconn, "menu_ramen", bson.M{})
	if err != nil {
		resp.Status = "error"
		resp.Message = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}

	resp.Status = "success"
	resp.Message = "Data retrieved successfully"
	resp.Data = resto

	helper.WriteJSON(respw, http.StatusOK, resp)
}

func GetMenuByOutletID(respw http.ResponseWriter, req *http.Request) {
	outletID := req.URL.Query().Get("outlet_id")
	if outletID == "" {
		respondWithError(respw, http.StatusBadRequest, "Outlet ID harus disertakan")
		return
	}

	objID, err := primitive.ObjectIDFromHex(outletID)
	if err != nil {
		respondWithError(respw, http.StatusBadRequest, "Outlet ID tidak valid")
		return
	}

	filter := bson.M{"outlet_id": objID}

	var menu []model.Menu
	menu, err = atdb.GetFilteredDocs[[]model.Menu](config.Mongoconn, "menu_ramen", filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(respw, http.StatusNotFound, "Menu tidak ditemukan untuk outlet ini")
		} else {
			respondWithError(respw, http.StatusInternalServerError, fmt.Sprintf("Terjadi kesalahan: %v", err))
		}
		return
	}

	// Return response JSON
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	json.NewEncoder(respw).Encode(map[string]interface{}{
		"status": "success",
		"data":   menu,
	})
}

// Helper function to respond with error
func respondWithError(respw http.ResponseWriter, code int, message string) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(code)
	json.NewEncoder(respw).Encode(map[string]string{"error": message})
}

func Postmenu_ramen(respw http.ResponseWriter, req *http.Request) {
	var restoran model.Menus

	// Decode request body
	if err := json.NewDecoder(req.Body).Decode(&restoran); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	// Validasi dan sanitasi input
	if restoran.NamaMenu == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "Nama menu tidak boleh kosong"})
		return
	}
	restoran.NamaMenu = html.EscapeString(restoran.NamaMenu)

	// Validasi harga (harus angka positif)
	if restoran.Harga <= 0 {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "Harga harus lebih besar dari 0"})
		return
	}

	// Sanitasi deskripsi untuk mencegah XSS
	restoran.Deskripsi = html.EscapeString(restoran.Deskripsi)

	// Validasi kategori
	if restoran.Kategori == "" {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "Kategori tidak boleh kosong"})
		return
	}
	restoran.Kategori = html.EscapeString(restoran.Kategori)

	// Simpan data yang sudah disanitasi ke dalam database
	result, err := config.Mongoconn.Collection("menu_ramen").InsertOne(context.Background(), restoran)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	// Ambil Inserted ID
	insertedID := result.InsertedID.(primitive.ObjectID)
	restoran.ID = insertedID

	// Kembalikan response sukses
	response := map[string]interface{}{
		"message":     "Menu ramen berhasil ditambahkan",
		"inserted_id": insertedID.Hex(),
		"data":        restoran,
	}

	helper.WriteJSON(respw, http.StatusOK, response)
}

func PutMenuflutter(respw http.ResponseWriter, req *http.Request, id string) {
	// Convert the ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var newMenu model.Menus
	// Decode the request body into the newMenu struct
	if err := json.NewDecoder(req.Body).Decode(&newMenu); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Decoded document:", newMenu)

	// Set the ID from the URL
	newMenu.ID = objectID

	filter := bson.M{"_id": newMenu.ID}
	updateFields := bson.M{
		"nama_menu": newMenu.NamaMenu,
		"harga":     newMenu.Harga,
		"deskripsi": newMenu.Deskripsi,
		"gambar":    newMenu.Gambar,
		"kategori":  newMenu.Kategori,
	}

	fmt.Println("Filter:", filter)
	fmt.Println("Update:", updateFields)

	result, err := atdb.UpdateOneDoc(config.Mongoconn, "menu_ramen", filter, updateFields)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, "Document not found or not modified")
		return
	}

	// Return the updated menu data
	helper.WriteJSON(respw, http.StatusOK, newMenu)
}

func PutMenu(respw http.ResponseWriter, req *http.Request) {
	var newMenu model.Menus
	// Decode the request body into the newMenu struct
	if err := json.NewDecoder(req.Body).Decode(&newMenu); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Decoded document:", newMenu)

	// Jika ID kosong, tidak bisa melanjutkan
	if newMenu.ID.IsZero() {
		helper.WriteJSON(respw, http.StatusBadRequest, "ID is missing or invalid")
		return
	}

	filter := bson.M{"_id": newMenu.ID}
	updateFields := bson.M{
		"nama_menu": newMenu.NamaMenu,
		"harga":     newMenu.Harga,
		"deskripsi": newMenu.Deskripsi,
		"gambar":    newMenu.Gambar,
		"kategori":  newMenu.Kategori,
	}

	fmt.Println("Filter:", filter)
	fmt.Println("Update:", updateFields)

	result, err := atdb.UpdateOneDoc(config.Mongoconn, "menu_ramen", filter, updateFields)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, "Document not found or not modified")
		return
	}

	// Mengembalikan data menu yang sudah diperbarui
	helper.WriteJSON(respw, http.StatusOK, newMenu)
}

func DeleteMenu(respw http.ResponseWriter, req *http.Request) {
	var requestBody struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid JSON data"})
		return
	}

	// Convert ID to ObjectID
	objectId, err := primitive.ObjectIDFromHex(requestBody.ID)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid ID format"})
		return
	}

	// Create filter
	filter := bson.M{"_id": objectId}

	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "menu_ramen", filter)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to delete document", "error": err.Error()})
		return
	}

	if deleteResult.DeletedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "Document not found"})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}

func DeleteMenuflutter(respw http.ResponseWriter, req *http.Request, id string) {
	// Konversi ID dari string ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid ID format",
			"id":      id,
		})
		return
	}

	// Membuat filter berdasarkan ID
	filter := bson.M{"_id": objectID}

	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "menu_ramen", filter)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to delete document",
			"error":   err.Error(),
			"id":      id,
		})
		return
	}

	// Jika tidak ada dokumen yang dihapus
	if deleteResult.DeletedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"message": "Document not found",
			"id":      id,
		})
		return
	}

	// Jika berhasil dihapus
	helper.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":       "success",
		"message":      "Document deleted successfully",
		"deletedCount": deleteResult.DeletedCount,
		"deletedId":    id,
	})
}

func GetPesanan(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	orders, err := atdb.GetAllDoc[[]model.Pesanan](config.Mongoconn, "pesanan", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, orders)
}

func GetPesananByOutletID(respw http.ResponseWriter, req *http.Request) {
	// Ambil query parameter outlet_id
	outletID := req.URL.Query().Get("outlet_id")
	if outletID == "" {
		respondWithError(respw, http.StatusBadRequest, "Outlet ID harus disertakan")
		return
	}

	// Konversi outlet_id menjadi ObjectID MongoDB
	objID, err := primitive.ObjectIDFromHex(outletID)
	if err != nil {
		respondWithError(respw, http.StatusBadRequest, "Outlet ID tidak valid")
		return
	}

	// Filter berdasarkan outlet_id
	filter := bson.M{"outlet_id": objID}

	// Ambil data pesanan dari koleksi
	var pesanan []model.Pesanan
	pesanan, err = atdb.GetFilteredDocs[[]model.Pesanan](config.Mongoconn, "pesanan", filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondWithError(respw, http.StatusNotFound, "Pesanan tidak ditemukan untuk outlet ini")
		} else {
			respondWithError(respw, http.StatusInternalServerError, fmt.Sprintf("Terjadi kesalahan: %v", err))
		}
		return
	}

	// Return response JSON
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	json.NewEncoder(respw).Encode(map[string]interface{}{
		"status": "success",
		"data":   pesanan,
	})
}

func isValidObjectID(id string) bool {
	if len(id) != 24 {
		return false
	}
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// Fungsi handler untuk memproses ID
func GetPesananByID(respw http.ResponseWriter, req *http.Request) {
	pesananID := req.URL.Query().Get("id")
	if pesananID == "" {
		respondWithError(respw, http.StatusBadRequest, "Pesanan ID harus disertakan")
		return
	}

	// Validasi ID apakah valid ObjectID
	if !isValidObjectID(pesananID) {
		respondWithError(respw, http.StatusBadRequest, "Pesanan ID tidak valid")
		return
	}

	// Konversi ID menjadi ObjectID MongoDB
	objID, err := primitive.ObjectIDFromHex(pesananID)
	if err != nil {
		respondWithError(respw, http.StatusBadRequest, "Pesanan ID tidak valid")
		return
	}

	// Filter berdasarkan ID
	filter := bson.M{"_id": objID}
	var pesanan []model.Pesanan
	pesanan, err = atdb.GetFilteredDocs[[]model.Pesanan](config.Mongoconn, "pesanan", filter, nil)
	if err != nil || len(pesanan) == 0 {
		if err == mongo.ErrNoDocuments || len(pesanan) == 0 {
			respondWithError(respw, http.StatusNotFound, "Pesanan tidak ditemukan")
		} else {
			respondWithError(respw, http.StatusInternalServerError, fmt.Sprintf("Terjadi kesalahan: %v", err))
		}
		return
	}

	// Response data pesanan
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	json.NewEncoder(respw).Encode(map[string]interface{}{
		"status": "success",
		"data":   pesanan[0],
	})
}
func GetMenuByID(respw http.ResponseWriter, req *http.Request) {
	menuID := req.URL.Query().Get("id")
	if menuID == "" {
		respondWithError(respw, http.StatusBadRequest, "Pesanan ID harus disertakan")
		return
	}

	// Validasi ID apakah valid ObjectID
	if !isValidObjectID(menuID) {
		respondWithError(respw, http.StatusBadRequest, "Pesanan ID tidak valid")
		return
	}

	// Konversi ID menjadi ObjectID MongoDB
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		respondWithError(respw, http.StatusBadRequest, "Menu ID tidak valid")
		return
	}

	// Filter berdasarkan ID
	filter := bson.M{"_id": objID}
	var menu []model.Menus
	menu, err = atdb.GetFilteredDocs[[]model.Menus](config.Mongoconn, "menu_ramen", filter, nil)
	if err != nil || len(menu) == 0 {
		if err == mongo.ErrNoDocuments || len(menu) == 0 {
			respondWithError(respw, http.StatusNotFound, "Menu tidak ditemukan")
		} else {
			respondWithError(respw, http.StatusInternalServerError, fmt.Sprintf("Terjadi kesalahan: %v", err))
		}
		return
	}

	// Response data pesanan
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	json.NewEncoder(respw).Encode(map[string]interface{}{
		"status": "success",
		"data":   menu[0],
	})
}

func GetPesananByStatus(respw http.ResponseWriter, req *http.Request) {

	status := req.URL.Query().Get("status")
	if status == "" {
		http.Error(respw, "Status pesanan harus disertakan", http.StatusBadRequest)
		return
	}

	validStatuses := []string{"baru", "diproses", "selesai"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		http.Error(respw, "Status pesanan tidak valid", http.StatusBadRequest)
		return
	}

	fmt.Println("Received status:", status)

	filter := bson.M{"status_pesanan": status}

	fmt.Println("Filter untuk MongoDB:", filter)

	var pesanan []model.Pesanan
	pesanan, err := atdb.GetFilteredDocs[[]model.Pesanan](config.Mongoconn, "pesanan", filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(respw, "Pesanan tidak ditemukan dengan status ini", http.StatusNotFound)
		} else {
			http.Error(respw, fmt.Sprintf("Terjadi kesalahan: %v", err), http.StatusInternalServerError)
		}
		return
	}

	helper.WriteJSON(respw, http.StatusOK, pesanan)
}

func GetPesananByStatusflutter(respw http.ResponseWriter, req *http.Request) {
	var resp struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	status := req.URL.Query().Get("status")
	if status == "" {
		resp.Status = "error"
		resp.Message = "Status pesanan harus disertakan"
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}

	validStatuses := []string{"baru", "diproses", "selesai"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		resp.Status = "error"
		resp.Message = "Status pesanan tidak valid"
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}

	filter := bson.M{"status_pesanan": status}

	var pesanan []model.Pesanan
	pesanan, err := atdb.GetFilteredDocs[[]model.Pesanan](config.Mongoconn, "pesanan", filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp.Status = "error"
			resp.Message = "Pesanan tidak ditemukan dengan status ini"
			helper.WriteJSON(respw, http.StatusNotFound, resp)
		} else {
			resp.Status = "error"
			resp.Message = fmt.Sprintf("Terjadi kesalahan: %v", err)
			helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		}
		return
	}

	resp.Status = "success"
	resp.Message = "Data retrieved successfully"
	resp.Data = pesanan

	helper.WriteJSON(respw, http.StatusOK, resp)
}

func PostPesanan(respw http.ResponseWriter, req *http.Request) {
	var pesanan model.Pesanan

	// Decode request body into pesanan
	if err := json.NewDecoder(req.Body).Decode(&pesanan); err != nil {
		log.Println("Error decoding request body:", err)
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "Invalid request body"})
		return
	}

	if len(pesanan.DaftarMenu) == 0 {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: "Daftar menu cannot be empty"})
		return
	}

	pesanan.StatusPesanan = "baru"
	pesanan.Pembayaran = "Cash"
	pesanan.TanggalPesanan = primitive.NewDateTimeFromTime(time.Now())

	log.Println("Pesanan data received:", pesanan)

	result, err := config.Mongoconn.Collection("pesanan").InsertOne(context.Background(), pesanan)
	if err != nil {
		log.Println("Error inserting pesanan:", err)
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: "Failed to insert pesanan"})
		return
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Error asserting inserted ID to ObjectID")
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: "Failed to process inserted ID"})
		return
	}

	log.Println("Pesanan berhasil disimpan dengan ID:", insertedID.Hex())

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Pesanan berhasil disimpan dengan ID: %s", insertedID.Hex())})
}

func GetItemPesanan(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	items, err := atdb.GetAllDoc[[]model.ItemPesanan](config.Mongoconn, "item_pesanan", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, items)
}

func PostItemPesanan(respw http.ResponseWriter, req *http.Request) {
	var item model.ItemPesanan
	if err := json.NewDecoder(req.Body).Decode(&item); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	result, err := config.Mongoconn.Collection("item_pesanan").InsertOne(context.Background(), item)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)
	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Item pesanan berhasil disimpan dengan ID: %s", insertedID.Hex())})
}

func CompleteOrder(respw http.ResponseWriter, req *http.Request) {

	orderID := req.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(respw, "Order ID harus disertakan", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		http.Error(respw, "Order ID tidak valid", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"status_pesanan": "Selesai",
			"waktu_terima":   primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := config.Mongoconn.Collection("pesanan").UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(respw, fmt.Sprintf("Terjadi kesalahan: %v", err), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(respw, "Pesanan tidak ditemukan", http.StatusNotFound)
		return
	}

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: "Pesanan berhasil diselesaikan"})
}

func UpdatePesananStatus(respw http.ResponseWriter, req *http.Request) {
	respw.Header().Set("Content-Type", "application/json")

	pesananID := req.URL.Query().Get("id")
	if pesananID == "" {
		http.Error(respw, `{"error": "ID pesanan harus disertakan"}`, http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(pesananID)
	if err != nil {
		http.Error(respw, `{"error": "ID pesanan tidak valid"}`, http.StatusBadRequest)
		return
	}

	var requestBody struct {
		StatusPesanan string `json:"status_pesanan"`
	}
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		http.Error(respw, `{"error": "Gagal membaca body request"}`, http.StatusBadRequest)
		return
	}
	defer req.Body.Close() // Hindari kebocoran resource

	validStatuses := map[string]bool{"baru": true, "diproses": true, "selesai": true}
	if !validStatuses[requestBody.StatusPesanan] {
		http.Error(respw, `{"error": "Status pesanan tidak valid"}`, http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status_pesanan": requestBody.StatusPesanan}}

	result, err := config.Mongoconn.Collection("pesanan").UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(respw, fmt.Sprintf(`{"error": "Terjadi kesalahan pada server: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(respw, `{"error": "Pesanan tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	// 🔥 Perbaikan Utama: Gunakan json.Marshal() agar tidak ada newline tambahan
	response := map[string]string{"message": "Status pesanan berhasil diperbarui"}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(respw, `{"error": "Gagal mengencode respons"}`, http.StatusInternalServerError)
		return
	}

	// Pastikan tidak ada karakter tambahan yang dikirim
	respw.WriteHeader(http.StatusOK)
	respw.Write(jsonData) // ✅ Mengirim JSON langsung tanpa newline tambahan
}
