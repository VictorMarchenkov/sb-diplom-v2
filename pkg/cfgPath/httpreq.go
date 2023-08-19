package config

import (
	"errors"
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
		return errors.New("error path to data of sms undefined")
	}

	if c.Email == "" {
		return errors.New("error path to data of email undefined")
	}

	if c.Voice == "" {
		return errors.New("error path to data of voice undefined")
	}

	if c.Billing == "" {
		return errors.New("error path to data of billing undefined")
	}

	return nil
}
