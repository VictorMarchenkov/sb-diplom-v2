package email

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sb-diplom-v2/internal/entities/consts"
)

type Data struct {
	country      string
	provider     string
	deliveryTime int
}

func (d *Data) validate() error {
	if len(d.country) != 2 {
		return fmt.Errorf("invalid country code length: %v", d.country)
	}

	if _, ok := consts.ISOCountries[d.country]; !ok {
		return fmt.Errorf("unknown country code: %v", d.country)
	}

	if d.provider == "" {
		return fmt.Errorf("provider is empty")
	}

	if _, ok := consts.EmailProviders[d.provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.provider)
	}

	return nil
}

type Records []Data

func newFromString(str string) (Data, error) {
	fields := strings.Split(str, ";")
	if len(fields) != 3 {
		return Data{}, fmt.Errorf("wrong number of fields: %d", len(fields))
	}
	result := Data{
		country:  fields[1],
		provider: fields[2],
	}

	v, err := strconv.Atoi(fields[3])
	if err != nil {
		return Data{}, fmt.Errorf("invalid ttfb: %v", fields[3])
	}
	result.deliveryTime = v

	if err := result.validate(); err != nil {
		return Data{}, err
	}

	return result, nil
}

func NewFromFile(fileName string) (Records, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result Records
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
