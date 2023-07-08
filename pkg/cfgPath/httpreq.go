package cfgPath

import (
	"fmt"
)

type httpService struct {
	Port     int    `json:"port"`
	Status   string `json:"status"`
	Mms      string `json:"mms"`
	Support  string `json:"support"`
	Incident string `json:"incident"`
}

func (s *httpService) initiate() error {
	return s.validate()
}

func (s *httpService) validate() error {

	if s.Port == 0 {
		return fmt.Errorf("undefined service port")
	}

	if s.Port > 65535 {
		return fmt.Errorf("invalide service port")
	}

	if s.Mms == "" {
		return fmt.Errorf("path to data of Mms info undefined")
	}

	if s.Support == "" {
		return fmt.Errorf("path to data of Support info undefined")
	}

	if s.Incident == "" {
		return fmt.Errorf("path to data of Incident info undefined")
	}

	return nil
}
