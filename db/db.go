package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq" //import postgres
	"github.com/restsec/api-gingonic/config"
)

//DB ...
type DB struct {
	*sql.DB
}

var db *gorp.DbMap

//Init ...
func Init() {

	dbinfo := fmt.Sprintf("host= %s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.ConfigParams.DbHost, config.ConfigParams.DbPort, config.ConfigParams.DbUser,
		config.ConfigParams.DbPassword, config.ConfigParams.DbName)

	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}

}

//ConnectDB ...
func ConnectDB(dataSourceName string) (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.TraceOn("[DB-GIN]", log.New(config.LogFile, "golang-GIN:", log.Lmicroseconds)) //Trace database requests
	return dbmap, nil
}

//GetDB ...
func GetDB() *gorp.DbMap {
	return db
}
