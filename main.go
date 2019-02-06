package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/auyer/muxapi/config"
	"github.com/auyer/muxapi/controllers"
	"github.com/auyer/muxapi/db"

	"github.com/gorilla/mux"
)

func main() {
	longConfigFile := flag.String("config", "", "Use to indicate the configuration file location")
	shortConfigFile := flag.String("c", "", "Use to indicate the configuration file location")
	flag.Parse()
	var conf config.ConfigurationStruct
	var err error
	if *longConfigFile != "" || *shortConfigFile != "" {
		conf, err = config.ReadFromFile(*longConfigFile + *shortConfigFile)
		if err != nil {
			log.Print(err.Error())
			return
		}
	} else {
		conf = config.ReadFromEnv()
	}

	fmt.Println("Starting GORILLA MUX API")

	log.SetOutput(conf.LogFile)
	dbPointer, err := db.ConnectDB(conf.DbAddress)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer dbPointer.Close()
	controller := new(controllers.Controller)
	controller.DB = dbPointer

	httpsRouter := mux.NewRouter()

	httpsRouter.HandleFunc("/api/", controller.GetServidor).Methods("GET")
	httpsRouter.HandleFunc("/api/", controller.PostServidor).Methods("POST")
	httpsRouter.HandleFunc("/api/{id:[0-9]+}", controller.GetServidorMat).Methods("GET") // URL parameter with Regex in URL
	httpsRouter.HandleFunc("/api/websocket/", controller.WebsocketHandler).Methods("GET")
	err = http.ListenAndServeTLS(":"+conf.HTTPSPort, conf.TLSCertLocation, conf.TLSKeyLocation, httpsRouter)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}
}
