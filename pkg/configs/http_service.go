package configs

import "errors"

type httpService struct {
	MMSURL      string `json:"mms_url"`
	SupportURL  string `json:"support_url"`
	IncidentURL string `json:"incident_url"`
	StatusURL   string `json:"status_url"`
	Port        int    `json:"port"`
	Status      string `json:"status"`
	Mms         string `json:"mms"`
	Support     string `json:"support"`
	Incident    string `json:"incident"`
}

func (s *httpService) initiate() error {
	return s.validate()
}

func (s *httpService) validate() error {
	if s.MMSURL == "" {
		return errors.New("error empty mms_url")
	}

	if s.IncidentURL == "" {
		return errors.New("error empty incident_url")
	}

	if s.SupportURL == "" {
		return errors.New("error empty support_url")
	}

	if s.StatusURL == "" {
		return errors.New("error empty status_url")
	}

	if s.Port == 0 {
		return errors.New("error empty port")
	}

	if s.Port < 0 || s.Port > 655535 {
		return errors.New("error invalid port value, must be from 1 to 65534")
	}

	if s.Mms == "" {
		return errors.New("error path to data of Mms info undefined")
	}

	if s.Support == "" {
		return errors.New("error path to data of Support info undefined")
	}

	if s.Incident == "" {
		return errors.New("error path to data of Incident info undefined")
	}

	if s.Status == "" {
		return errors.New("error path to data of Status undefined")
	}

	return nil
}
