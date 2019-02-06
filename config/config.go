package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// ConfigurationStruct models the necessary configs for an API
type ConfigurationStruct struct {
	LogLocation     string `json:"LogLocation"`
	HTTPPort        string `json:"HTTPPort"`
	HTTPSPort       string `json:"HTTPSPort"`
	TLSKeyLocation  string `json:"TLSKeyLocation"`
	TLSCertLocation string `json:"TLSCertLocation"`
	DbAddress       string `json:"DbAddress"`
	Debug           string `json:"Debug"`
	LogFile         io.Writer
}

// ReadFromEnv reads the configuration from enviroment variables
func ReadFromEnv() ConfigurationStruct {
	var configParams ConfigurationStruct
	configParams.LogLocation = os.Getenv("LOGLOCATION")
	configParams.LogFile = setLogFile(configParams.LogLocation)
	configParams.HTTPPort = os.Getenv("HTTPPORT")
	configParams.HTTPSPort = os.Getenv("HTTPSPORT")
	configParams.TLSKeyLocation = os.Getenv("TLSKEYLOCATION")
	configParams.TLSCertLocation = os.Getenv("TLSCERTLOCATION")
	configParams.DbAddress = os.Getenv("DBADDRESS")
	return configParams

}

// ReadFromFile reads from a config file or enviroment variables
func ReadFromFile(fileloc string) (ConfigurationStruct, error) {
	var configParams ConfigurationStruct
	file, err := ioutil.ReadFile(fileloc)

	if err != nil {
		log.Print(err.Error())
		return configParams, err
	}

	err = json.Unmarshal(file, &configParams)

	if err != nil {
		log.Print(err.Error())
		return configParams, err
	}
	if configParams.Debug == "true" {
		fmt.Println("Configuration Parameters")
		fmt.Println(string(file))
	}

	configParams.LogFile = setLogFile(configParams.LogLocation)

	return configParams, nil
}

func setLogFile(logLocation string) io.Writer {
	if logLocation != "" {
		if _, err := os.Stat(logLocation); os.IsNotExist(err) {
			fileLog, fileErr := os.Create(logLocation)
			if fileErr == nil {
				fmt.Println("Writing logs to file " + logLocation)
				return fileLog
			}
			fmt.Println(fileErr)
		} else {
			fileLog, fileErr := os.OpenFile(logLocation, os.O_RDWR|os.O_APPEND, 0660)
			if fileErr == nil {
				fmt.Println("Writing logs to file " + logLocation)
				return fileLog
			}
			fmt.Println(fileErr)
		}
	}
	return os.Stdout

}
