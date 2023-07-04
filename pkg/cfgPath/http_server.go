package cfgPath

import "fmt"

type httpServer struct {
	Port uint `json:"port"`
}

func (h *httpServer) initiate() error {
	return h.validate()
}

func (h *httpServer) validate() error {
	if h.Port == 0 {
		return fmt.Errorf("port is empty")
	}

	if h.Port > 65535 {
		return fmt.Errorf("invalid port")
	}

	return nil
}
