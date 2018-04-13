package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	LogFile      io.Writer
	ConfigParams *configStruct
)

type configStruct struct {
	LogLocation     string `json:"LogLocation"`
	HttpPort        string `json:"HttpPort"`
	HttpsPort       string `json:"HttpsPort"`
	TLSKeyLocation  string `json:"TLSKeyLocation"`
	TLSCertLocation string `json:"TLSCertLocation"`
	DbUser          string `json:"DatabaseUser"`
	DbPassword      string `json:"DatabasePassword"`
	DbName          string `json:"DatabaseName"`
	DbHost          string `json:"DatabaseHost"`
	DbPort          string `json:"DatabasePort"`
	Debug           string `json:"Debug"`
}

func ReadConfig() error {

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = json.Unmarshal(file, &ConfigParams)

	if err != nil {
		log.Print(err.Error())
		return err
	}
	if ConfigParams.Debug == "true" {
		fmt.Println("Configuration Parameters")
		fmt.Println(string(file))
	}

	if ConfigParams.LogLocation != "" {
		if _, err := os.Stat(ConfigParams.LogLocation); os.IsNotExist(err) {
			fileLog, fileErr := os.Create(ConfigParams.LogLocation)
			if fileErr != nil {
				fmt.Println(fileErr)
				fileLog = os.Stdout
			} else {
				fmt.Println("Writing logs to file " + ConfigParams.LogLocation)
			}
			LogFile = fileLog
		} else {
			fileLog, fileErr := os.OpenFile(ConfigParams.LogLocation, os.O_RDWR|os.O_APPEND, 0660)
			if fileErr != nil {
				fmt.Println(fileErr)
				fileLog = os.Stdout
			} else {
				fmt.Println("Writing logs to file " + ConfigParams.LogLocation)
			}
			LogFile = fileLog
		}
	} else {
		LogFile = os.Stdout
	}
	return nil
}
