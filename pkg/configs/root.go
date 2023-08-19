package configs

import (
	"errors"
)

type Root struct {
	HTTPServer  httpServer  `json:"http_server"`
	CSV         csv         `json:"csv"`
	HTTPService httpService `json:"http_service"`
}

func (r *Root) initiate() error {

	if err := r.HTTPServer.initiate(); err != nil {
		return errors.New("error initiate http_server")
	}

	if err := r.CSV.initiate(); err != nil {
		return errors.New("error initiate csv service")
	}

	if err := r.HTTPService.initiate(); err != nil {
		return errors.New("error initiate http service")
	}

	return nil
}
