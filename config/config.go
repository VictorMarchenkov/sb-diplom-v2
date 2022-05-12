package config

import (
	"io/ioutil"
	"log"
)

// GetConfig reading config file.
func GetConfig() []byte {
	rawConfig, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Println(err.Error())
		//log.Fatal(err.Error())
	}
	return rawConfig
}
