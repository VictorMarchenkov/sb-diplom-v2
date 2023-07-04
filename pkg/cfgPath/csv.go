package cfgPath

import (
	"fmt"
)

type csv struct {
	Sms     string `json:"sms"`
	Voice   string `json:"voice"`
	Email   string `json:"email"`
	Billing string `json:"billing"`
}

func (c *csv) initiate() error {
	return c.validate()
}

func (c *csv) validate() error {
	if c.Sms == "" {
		return fmt.Errorf("path to data of sms undefined")
	}

	if c.Email == "" {
		return fmt.Errorf("path to data of email undefined")
	}

	if c.Voice == "" {
		return fmt.Errorf("path to data of voice undefined")
	}

	if c.Billing == "" {
		return fmt.Errorf("path to data of billing undefined")
	}

	return nil
}
