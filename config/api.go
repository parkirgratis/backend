package config

import "os"

var WAAPIQRLogin string = "https://api.wa.my.id/api/whatsauth/request"

var WAAPIMessage string = "https://api.wa.my.id/api/send/message/text"

var WAAPIGetToken string = "https://api.wa.my.id/api/signup"

var PublicKeyWhatsAuth string

var WAAPITextMessage string = "https://api.wa.my.id/api/v2/send/message/text"

var WAAPIDocMessage string = "https://api.wa.my.id/api/send/message/document"

var APITOKENPD string = os.Getenv("PDTOKEN")

var WAAPIToken string

var APIGETPDLMS string = "https://pamongdesa.kemendagri.go.id/webservice/public/user/get-by-phone?number="

var WAAPIGetDevice string = "https://api.wa.my.id/api/device/"

var GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail string

