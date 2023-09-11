package configs

import (
	"encoding/json"
	"errors"
	"os"
)

func Init(filePath string) (*Root, error) {

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("error read file")
	}

	var root Root
	if err := json.Unmarshal(f, &root); err != nil {
		return nil, errors.New("error parse config")
	}

	if err := root.initiate(); err != nil {
		return nil, err
	}

	return &root, nil
}
