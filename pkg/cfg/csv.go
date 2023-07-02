package cfg

import (
	"fmt"
	entities "sb-diplom-v2/internal"
)

type csv struct {
	entities.Csv
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
