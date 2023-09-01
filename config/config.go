package config

import (
	"errors"
	"io/ioutil"
)

// GetConfig reading config file.
func GetConfig() ([]byte, error) {
	rawConfig, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		errors.New("error opening config file")
		return nil, err
	}
	return rawConfig, nil
}
