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
	case method == "GET" && path == "/data/lokasi":
		controller.GetLokasi(w, r)
	case method == "GET" && path == "/data/marker":
		controller.GetMarker(w, r)
	case method == "GET" && path == "/data/nama-tempat":
		controller.GetTempatByNamaTempat(w, r)
	case method == "POST" && helper.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)
	case method == "POST" && path == "/tempat-parkir":
		controller.PostTempatParkir(w, r)
	case method == "POST" && path == "/koordinat":
		controller.PostKoordinat(w, r)
	case method == "POST" && helper.URLParam(path, "/upload/:path"):
		controller.PostUploadGithub(w, r)
	case method == "PUT" && path == "/data/tempat":
		controller.PutTempatParkir(w, r)
	case method == "PUT" && path == "/data/koordinat":
		controller.PutKoordinat(w, r)
	case method == "DELETE" && path == "/data/tempat":
		controller.DeleteTempatParkir(w, r)
	case method == "DELETE" && path == "/data/koordinat":
		controller.DeleteKoordinat(w, r)
	case method == "POST" && path == "/admin/login":
		handler.Login(w, r)
	case method == "POST" && path == "/admin/logout":
		handler.Logout(w, r)
	case method == "GET" && path == "/data/saran":
		controller.GetSaran(w, r)
	case method == "POST" && path == "/data/saran":
		controller.PostSaran(w, r)
	case method == "PUT" && path == "/data/saran":
		controller.PutSaran(w, r)
	case method == "DELETE" && path == "/data/saran":
		controller.DeleteSaran(w, r)
	case method == "GET" && path == "/admin/dashboard":
		middleware.AuthMiddleware(http.HandlerFunc(handler.DashboardAdmin)).ServeHTTP(w, r)
	default:
		controller.NotFound(w, r)
	}
}
