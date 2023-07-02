package cfg

import (
	"fmt"
	entities "sb-diplom-v2/internal"
)

type httpReq struct {
	entities.Http
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
