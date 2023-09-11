package configs

import "errors"

type httpService struct {
	MMSURL      string `json:"mms_url"`
	SupportURL  string `json:"support_url"`
	IncidentURL string `json:"incident_url"`
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

	return nil
}
