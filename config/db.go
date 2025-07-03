package config

import (
	"os"

	"github.com/gocroot/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")

// Koneksi ke database parkir_db
var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "parkir_db",
}
var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)

// Koneksi ke database ramen
var ramenInfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "ramen",
}
var RamenConn, ErrorRamenConn = atdb.MongoConnect(ramenInfo)
