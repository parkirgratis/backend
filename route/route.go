package route

import (
	"net/http"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/controller"
	"github.com/gocroot/handler"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/middleware"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return
	}
	config.SetEnv()

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		controller.GetHome(w, r)

		// Rute untuk fitur ParkirGratis
	case method == "GET" && path == "/data/lokasi":
		controller.GetLokasi(w, r) // Mengambil data lokasi parkir.
	case method == "GET" && path == "/data/marker":
		controller.GetMarker(w, r) // Mendapatkan data marker lokasi.
	case method == "GET" && path == "/data/search-namatempat":
		controller.GetTempatByNamaTempat(w, r) // Pencarian berdasarkan nama tempat.
		// Rute untuk mengelola data tempat parkir.
	case method == "POST" && path == "/tempat-parkir":
		controller.PostTempatParkir(w, r) // Menambahkan tempat parkir.
	case method == "PUT" && path == "/data/tempat":
		controller.PutTempatParkir(w, r) // Memperbarui data tempat parkir.
	case method == "DELETE" && path == "/data/tempat":
		controller.DeleteTempatParkir(w, r) // Menghapus tempat parkir.
	case method == "POST" && path == "/upload/img":
		controller.PostUploadGithub(w, r)

	case method == "POST" && path == "/tempat-parkir":
		controller.InsertDataRegionFromPetapdia(w, r)

		// Rute untuk mengelola koordinat.
	case method == "POST" && path == "/koordinat":
		controller.PostKoordinat(w, r) // Menambahkan koordinat baru.

		// Rute untuk mengelola saran.
	case method == "GET" && path == "/data/saran":
		controller.GetSaran(w, r) // Mendapatkan daftar saran.
	case method == "POST" && path == "/data/datasaran":
		controller.PostSaran(w, r) // Menambahkan saran baru.
	case method == "PUT" && path == "/data/saran":
		controller.PutSaran(w, r) // Memperbarui data saran.
	case method == "DELETE" && path == "/data/saran":
		controller.DeleteSaran(w, r) // Menghapus data saran.

	// Rute untuk admin (login, logout, register, dashboard, aktivitas).
	case method == "POST" && path == "/admin/login":
		handler.Login(w, r) // Login admin.
	case method == "POST" && path == "/admin/logout":
		handler.Logout(w, r) // Logout admin.
	case method == "POST" && path == "/admin/register":
		handler.RegisterAdmin(w, r) // Registrasi admin baru.
	case method == "GET" && path == "/admin/dashboard":
		// Middleware autentikasi untuk dashboard admin.
		middleware.AuthMiddleware(http.HandlerFunc(handler.DashboardAdmin)).ServeHTTP(w, r)

		// Rute untuk fitur Warung.
	case method == "POST" && path == "/data/tempat-warung":
		controller.PostTempatWarung(w, r) // Menambahkan data warung.
	case method == "GET" && path == "/data/warung":
		controller.GetaAllWarung(w, r) // Mendapatkan semua data warung.
	case method == "DELETE" && path == "/data/deletewarung":
		controller.DeleteTempatWarungById(w, r) // Delete data warung berdasarkan Id.
	case method == "PUT" && path == "/data/warung":
		controller.PutTempatWarung(w, r) // Update/Edit data warung berdasarkan Id.
	case method == "GET" && path == "/data/markerwarung":
		controller.GetMarkerWarung(w, r)
	case method == "PUT" && path == "/data/markerwarung":
		controller.PutKoordinatWarung(w, r)
	case method == "DELETE" && path == "/data/markerwarung":
		controller.DeleteKoordinatWarung(w, r)
	case method == "POST" && path == "/data/marker-warung":
		controller.PostKoordinatWarung(w, r)
	case method == "GET" && path == "/data/search-namatempatwarung":
		controller.GetByNamaTempatWarung(w, r)

		//Location Nembak Endpoint Dari Petapedia
	case method == "POST" && path == "/data/gis/lokasiparkir":
		controller.InsertDataRegionFromPetapdia(w, r)
	case method == "POST" && path == "/data/gis/lokasiwarung":
		controller.InsertDataRegionFromPetapdiaWarung(w, r)
	case method == "POST" && path == "/data/gis/jalan":
		controller.GetRoads(w, r)
	case method == "POST" && path == "/data/gis/region":
		controller.GetRegion(w, r)

	// Rute untuk webhook dengan parameter dinamis.
	case method == "POST" && helper.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)

	// Google Auth
	case method == "POST" && path == "/auth/users":
		controller.Auth(w, r)
	case method == "POST" && path == "/auth/login":
		controller.GeneratePasswordHandler(w, r)
	case method == "POST" && path == "/auth/verify":
		controller.VerifyPasswordHandler(w, r)
	case method == "POST" && path == "/auth/resend":
		controller.ResendPasswordHandler(w, r)

	//user data
	case method == "GET" && path == "/data/user":
		controller.GetDataUser(w, r)
	//mendapatkan data sent item
	case method == "GET" && at.URLParam(path, "/data/peserta/sent/:id"):
		controller.GetSentItem(w, r)
	//simpan feedback unsubs user
	case method == "POST" && path == "/data/peserta/unsubscribe":
		controller.PostUnsubscribe(w, r)
	case method == "POST" && path == "/data/user":
		controller.PostDataUser(w, r)
	//generate token linked device
	case method == "PUT" && path == "/data/user":
		controller.PutTokenDataUser(w, r)
	//Menambhahkan data nomor sender untuk broadcast
	case method == "PUT" && path == "/data/sender":
		controller.PutNomorBlast(w, r)
	//mendapatkan data list nomor sender untuk broadcast
	case method == "GET" && path == "/data/sender":
		controller.GetDataSenders(w, r)
	//mendapatkan data list nomor sender yang kena blokir dari broadcast
	case method == "GET" && path == "/data/blokir":
		controller.GetDataSendersTerblokir(w, r)
	//mendapatkan data rekap pengiriman wa blast
	case method == "GET" && path == "/data/rekap":
		controller.GetRekapBlast(w, r)
	//mendapatkan data faq
	case method == "GET" && at.URLParam(path, "/data/faq/:id"):
		controller.GetFAQ(w, r)
	case method == "POST" && at.URLParam(path, "/data/user/wa/:nomorwa"):
		controller.PostDataUserFromWA(w, r)

	case method == "POST" && helper.URLParam(path, "/upload/:path"):
		controller.PostUploadGithub(w, r)

		//ramennn
		// endpoint menu ramen
	case method == "GET" && path == "/data/menu_ramen":
		controller.GetMenu_ramen(w, r)

	case method == "PUT" && path == "/ubah/menu_ramen":
		controller.PutMenu(w, r)

	case method == "GET" && path == "/menu/byid":
		controller.GetMenuByID(w, r)

	case method == "GET" && path == "/data/ramen":
		controller.GetMenu_ramenflutter(w, r)

	case method == "POST" && path == "/tambah/menu_ramen":
		controller.Postmenu_ramen(w, r)

	case method == "PUT" && strings.HasPrefix(path, "/ubah/byid/"):
		// Extract the ID from the path
		id := strings.TrimPrefix(path, "/ubah/byid/")
		// Call the PutMenu function with the extracted ID
		controller.PutMenuflutter(w, r, id)

	case method == "DELETE" && path == "/hapus/menu_ramen":
		controller.DeleteMenu(w, r)

	case method == "DELETE" && strings.HasPrefix(path, "/hapus/byid/"):
		// Ambil ID dari URL
		id := strings.TrimPrefix(path, "/hapus/byid/")
		// Panggil fungsi DeleteMenu dengan ID dari URL
		controller.DeleteMenuflutter(w, r, id)

		// endpoint pesanan
	case method == "GET" && path == "/data/pesanan":
		controller.GetPesanan(w, r)
	case method == "GET" && path == "/data/byid":
		controller.GetPesananByID(w, r)

	case method == "GET" && path == "/data/bystatus":
		controller.GetPesananByStatus(w, r)

	case method == "GET" && path == "/data/bystatus/flutter":
		controller.GetPesananByStatusflutter(w, r)

	case method == "POST" && path == "/tambah/pesanan":
		controller.PostPesanan(w, r)
	case method == "PATCH" && path == "/update/status":
		controller.UpdatePesananStatus(w, r)

		// endpoint item pesanan
		controller.GetItemPesanan(w, r)
	case method == "POST" && path == "/tambah/item_pesanan":
		controller.PostItemPesanan(w, r)
	case method == "POST" && helper.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)

		// Rute untuk admin (login, logout, register, dashboard, aktivitas).
	case method == "POST" && path == "/admin/logins":
		handler.Logins(w, r) // Login admin.
	case method == "GET" && path == "/data/activitys":
		controller.GetActivity(w, r)
	case method == "GET" && path == "/data/admin":
		handler.GetAllAdmins(w, r)
	case method == "POST" && path == "/admin/logouts":
		handler.Logouts(w, r) // Logout admin.
	case method == "POST" && path == "/admin/registers":
		handler.RegisterAdmins(w, r) // Registrasi admin baru.
		
	case method == "GET" && path == "/admin/dashboards":

	case method == "PUT" && path == "/update/password":
		handler.UpdateForgottenPassword(w, r) // Login admin.
		// Middleware autentikasi untuk dashboard admin.
		middleware.AuthMiddleware(http.HandlerFunc(handler.DashboardAdmins)).ServeHTTP(w, r)

	// Rute default untuk request yang tidak dikenali.
	default:
		controller.NotFound(w, r) // Mengembalikan response 404 jika rute tidak ditemukan.
	}
}
