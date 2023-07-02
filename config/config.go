package config

import (
	"fmt"
	"io/ioutil"
)

// GetConfig reading config file.
func GetConfig() ([]byte, error) {
	rawConfig, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("error opening config filr: ", err.Error())
		return nil, err
	}
	return rawConfig, nil
}
