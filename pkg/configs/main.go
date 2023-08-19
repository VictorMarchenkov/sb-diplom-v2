package configs

import (
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Println(root.HTTPServer.Port, root.HTTPService.Port)
	if err := root.initiate(); err != nil {
		return nil, err
	}

	return &root, nil
}
