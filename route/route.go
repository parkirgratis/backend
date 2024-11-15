package route

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/controller"
	"github.com/gocroot/handler"
	"github.com/gocroot/helper"
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
	case method == "GET" && path == "/data/user":
		controller.GetDataUser(w, r) // Mengambil data pengguna.
	case method == "PUT" && path == "/data/user":
		controller.PutTokenDataUser(w, r) // Update token pengguna.
	case method == "GET" && path == "/data/marker":
		controller.GetMarker(w, r) // Mendapatkan data marker lokasi.
	case method == "GET" && path == "/data/search-namatempat":
		controller.GetTempatByNamaTempat(w, r) // Pencarian berdasarkan nama tempat.

	// Rute untuk webhook dengan parameter dinamis.
	case method == "POST" && helper.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)

	// Rute untuk mengelola data tempat parkir.
	case method == "POST" && path == "/tempat-parkir":
		controller.PostTempatParkir(w, r) // Menambahkan tempat parkir.
	case method == "PUT" && path == "/data/tempat":
		controller.PutTempatParkir(w, r) // Memperbarui data tempat parkir.
	case method == "DELETE" && path == "/data/tempat":
		controller.DeleteTempatParkir(w, r) // Menghapus tempat parkir.

	// Rute untuk mengelola koordinat.
	case method == "POST" && path == "/koordinat":
		controller.PostKoordinat(w, r) // Menambahkan koordinat baru.
	case method == "PUT" && path == "/data/koordinat":
		controller.PutKoordinat(w, r) // Memperbarui data koordinat.
	case method == "DELETE" && path == "/data/koordinat":
		controller.DeleteKoordinat(w, r) // Menghapus koordinat.

	// Rute untuk admin (login, logout, register, dashboard, aktivitas).
	case method == "POST" && path == "/admin/login":
		handler.Login(w, r) // Login admin.
	case method == "POST" && path == "/admin/logout":
		handler.Logout(w, r) // Logout admin.
	case method == "POST" && path == "/admin/register":
		handler.RegisterAdmin(w, r) // Registrasi admin baru.
	case method == "POST" && path == "/admin/activity":
		controller.LogActivity(w, r) // Log aktivitas admin.
	case method == "GET" && path == "/admin/dashboard":
		// Middleware autentikasi untuk dashboard admin.
		middleware.AuthMiddleware(http.HandlerFunc(handler.DashboardAdmin)).ServeHTTP(w, r)

	// Rute untuk mengelola saran.
	case method == "GET" && path == "/data/saran":
		controller.GetSaran(w, r) // Mendapatkan daftar saran.
	case method == "POST" && path == "/data/saran":
		controller.PostSaran(w, r) // Menambahkan saran baru.
	case method == "PUT" && path == "/data/saran":
		controller.PutSaran(w, r) // Memperbarui data saran.
	case method == "DELETE" && path == "/data/saran":
		controller.DeleteSaran(w, r) // Menghapus data saran.

	// Rute untuk fitur Warung.
	case method == "POST" && path == "/data/warung":
		controller.PostTempatWarung(w, r) // Menambahkan data warung.
	case method == "GET" && path == "/data/warung":
		controller.GetaAllWarung(w, r) // Mendapatkan semua data warung.

	// Rute default untuk request yang tidak dikenali.
	default:
		controller.NotFound(w, r) // Mengembalikan response 404 jika rute tidak ditemukan.
	}
}
