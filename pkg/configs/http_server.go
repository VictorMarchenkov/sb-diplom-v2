package configs

import (
	"errors"
	"fmt"
)

type httpServer struct {
	Port uint `json:"port"`
}

func (h *httpServer) initiate() error {
	return h.validate()
}

func (h *httpServer) validate() error {
	if h.Port == 0 {
		return errors.New("port is empty")
	}

	if h.Port > 65535 {
		return errors.New("invalid port")
	}

	return nil
}

func (h *httpServer) HostPort() string {
	return fmt.Sprint(":", h.Port)
}
