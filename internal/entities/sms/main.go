package sms

import (
	"bufio"
	"fmt"
	"os"
<<<<<<< HEAD
	"sb-diplom-v2/internal/entities"
	"strconv"
	"strings"
)

type Data struct {
	country      string
	bandWidth    uint8
	responseTime uint
	provider     string
}

func (d *Data) validate() error {
	if len(d.country) != 2 {
		return fmt.Errorf("invalid country code length: %v", d.country)
	}

	if _, ok := entities.ISOCountries[d.country]; ok {
		return fmt.Errorf("unknown country code: %v", d.country)
	}

	if d.provider == "" {
		return fmt.Errorf("provider is empty")
	}

	if _, ok := entities.SMSProviders[d.provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.provider)
	}
	return nil
}

type Records []Data

func newFromString(str string) (Data, error) {
	fields := strings.Split(str, ";")
	if len(fields) != 4 {
		return Data{}, fmt.Errorf("wrong number of fields: %d", len(fields))
	}
	result := Data{
=======
	"strconv"
	"strings"

	"sb-diplom-v2/internal/entities"
)

type SMSData struct {
	country         string
	throughput      uint8
	avgResponseTime uint
	provider        string
}

func (s *SMSData) validate() error {
	if len(s.country) != 2 {
		return fmt.Errorf("invalid country code length: %v", s.country)
	}

	if _, ok := entities.ISOCountries[s.country]; !ok {
		return fmt.Errorf("unknown country code: %v", s.country)
	}

	if s.provider == "" {
		return fmt.Errorf("provider is empty")
	}

	if _, ok := entities.Providers[s.provider]; !ok {
		return fmt.Errorf("unknown provider: %v", s.provider)
	}

	return nil
}

type SMSRecords []SMSData

func newFromString(str string) (SMSData, error) {
	fields := strings.Split(str, ";")
	if len(fields) != 4 {
		return SMSData{}, fmt.Errorf("wrong number of fields: %d", len(fields))
	}

	result := SMSData{
>>>>>>> origin/master
		country:  fields[0],
		provider: fields[3],
	}

	v, err := strconv.ParseUint(fields[1], 10, 8)
	if err != nil {
<<<<<<< HEAD
		return Data{}, fmt.Errorf("invalid throughput: %q", fields[1])
	}
	result.bandWidth = uint8(v)

	v, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return Data{}, fmt.Errorf("invalid response time: %q", fields[2])
	}
	result.responseTime = uint(v)

	if err := result.validate(); err != nil {
		return Data{}, err
=======
		return SMSData{}, fmt.Errorf("invalid throughput: %q", fields[1])
	}
	result.throughput = uint8(v)

	v, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return SMSData{}, fmt.Errorf("invalid response time: %q", fields[2])
	}
	result.avgResponseTime = uint(v)

	if err := result.validate(); err != nil {
		return SMSData{}, err
>>>>>>> origin/master
	}

	return result, nil
}

<<<<<<< HEAD
func NewFromFile(fileName string) (Records, error) {
=======
func NewFromFile(fileName string) (SMSRecords, error) {
>>>>>>> origin/master
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

<<<<<<< HEAD
	var result Records
=======
	var result SMSRecords
>>>>>>> origin/master
	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if d, err := newFromString(line); err == nil {
			result = append(result, d)
		}
	}

	return result, nil
}
