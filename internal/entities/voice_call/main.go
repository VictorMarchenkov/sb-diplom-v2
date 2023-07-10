package voice_call

import (
	"bufio"
	"fmt"
	"os"
	"sb-diplom-v2/internal/entities"
	"strconv"
	"strings"
)

type Data struct {
	country             string
	bandwidth           string
	responseTime        string
	provider            string
	connectionStability float32
	ttfb                int
	voicePurity         int
	medianOfCallsTime   int
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

	if _, ok := entities.VoiceCallProviders[d.provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.provider)
	}

	return nil
}

type Records []Data

func newFromString(str string) (Data, error) {
	fields := strings.Split(str, ";")
	if len(fields) != 8 {
		return Data{}, fmt.Errorf("wrong number of fields: %d", len(fields))
	}

	result := Data{
		country:      fields[0],
		bandwidth:    fields[1],
		responseTime: fields[2],
		provider:     fields[3],
	}

	v4, err := strconv.ParseFloat(fields[4], 32)
	if err != nil {
		return Data{}, fmt.Errorf("invalid response time: %q", fields[4])
	}
	result.connectionStability = float32(v4)

	v5, err := strconv.Atoi(fields[5])
	if err != nil {
		return Data{}, fmt.Errorf("invalid ttfb: %v", fields[5])
	}
	result.ttfb = v5

	v6, err := strconv.Atoi(fields[6])
	if err != nil {
		return Data{}, fmt.Errorf("invalid ttfb: %v", fields[6])
	}
	result.voicePurity = v6

	v7, err := strconv.Atoi(fields[7])
	if err != nil {
		return Data{}, fmt.Errorf("invalid ttfb: %v", fields[7])
	}
	result.medianOfCallsTime = v7

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
