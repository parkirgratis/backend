package config

import (
	"net/http"
	"strings"
)

var AllowedOrigins = []string{
	"https://parkirgratis.github.io",
	"https://parkirgratis.github.io.id",
	"https://parkirgratis.if.co.id/input",
	"https://parkirgratis.if.co.id",
	"http://127.0.0.1:5500",
	"http://127.0.0.1:5501",
	"https://rayfanaqbil.github.io",
	"https://hrisz.github.io",
	"https://jul003.github.io",
	"https://farelnouval.github.io",
	"https://irgifauzi.github.io",
	"https://mhrndiva.github.io",
	"https://deviwlndr.github.io",
	"https://qinthargantheng.github.io",
	"https://serlip06.github.io",
	"http://agung6544.github.io",
	"https://yoginara.github.io",
	"https://barganakukuhraditya.github.io",
}

var AllowedHeaders = []string{
	"Origin",
	"Content-Type",
	"Accept",
	"Authorization",
	"Access-Control-Request-Headers",
	"Token",
	"Login",
	"Access-Control-Allow-Origin",
	"Bearer",
	"X-Requested-With",
}

func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	allowedOrigin := false
	for _, o := range AllowedOrigins {
		if o == origin {
			allowedOrigin = true
			break
		}
	}
	if !allowedOrigin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(AllowedHeaders, ", "))
	w.Header().Set("Access-Control-Allow-Origin", origin)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}
