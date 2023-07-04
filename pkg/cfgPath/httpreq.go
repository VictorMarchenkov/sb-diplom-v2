package cfgPath

import (
	"fmt"
)

type httpReq struct {
	ServerPort  int    `json:"server_port"`
	ServicePort int    `json:"service_port"`
	Status      string `json:"status"`
	Mms         string `json:"mms"`
	Support     string `json:"support"`
	Incident    string `json:"incident"`
}

func (h *httpReq) iniate() error {
	return h.validate()
}

func (h *httpReq) validate() error {
	if h.Mms == "" {
		return fmt.Errorf("path to data of Mms info undefined")
	}

	if h.Support == "" {
		return fmt.Errorf("path to data of Support info undefined")
	}

	if h.Incident == "" {
		return fmt.Errorf("path to data of Incident info undefined")
	}

	if h.ServerPort == 0 {
		return fmt.Errorf("undefined server port")
	}

	if h.ServerPort > 65535 {
		return fmt.Errorf("invalid server port")
	}

	if h.ServicePort == 0 {
		return fmt.Errorf("undefined service port")
	}

	if h.ServicePort > 65535 {
		return fmt.Errorf("invalide service port")
	}

	return nil
}
