package config

import (
	"os"

	"github.com/gocroot/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")
//mongostring
var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "parkir_db",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)