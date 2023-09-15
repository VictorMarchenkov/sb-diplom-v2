package voice_call

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"sb-diplom-v2/internal/entities/consts"
)

type Data struct {
	Country             string  `json:"country"`
	BandWidth           string  `json:"bandwidth"`
	ResponseTime        string  `json:"response_time"`
	Provider            string  `json:"provider"`
	ConnectionStability float32 `json:"connection_stability"`
	Ttfb                int     `json:"ttfb"`
	VoicePurity         int     `json:"voice_purity"`
	MedianOfCallsTime   int     `json:"median_of_calls_time"`
}

func (d *Data) validate() error {

	if len(d.Country) != 2 {
		return fmt.Errorf("invalid country code length: %v", d.Country)
	}

	if _, ok := consts.ISOCountries[d.Country]; !ok {
		return fmt.Errorf("unknown country code: %v", d.Country)
	}

	if d.Provider == "" {
		return fmt.Errorf("provider is empty")
	}

	if _, ok := consts.VoiceCallProviders[d.Provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.Provider)
	}

	if d.BandWidth == "" {
		return errors.New("band width undefined")
	}

	if d.ResponseTime == "" {
		return errors.New("response time width undefined")
	}

	var number int = 1
	var floatNumber float32 = 1

	if reflect.ValueOf(d.ConnectionStability).Type() != reflect.TypeOf(floatNumber) {
		return errors.New("wrong type for Connection Stability")
	}

	if reflect.ValueOf(d.MedianOfCallsTime).Type() != reflect.TypeOf(number) {
		return errors.New("wrong type for Median Of Calls Tim")
	}

	if reflect.ValueOf(d.Ttfb).Type() != reflect.TypeOf(number) {
		return errors.New("wrong type for ftfb")
	}

	if reflect.ValueOf(d.VoicePurity).Type() != reflect.TypeOf(number) {
		return errors.New("wrong type for Voice Purity")
	}

	return nil
}

type Set []Data

func decodeCSV(csvStr string) (Data, error) {

	fields := strings.Split(csvStr, ";")
	if len(fields) != 8 {
		return Data{}, fmt.Errorf("wrong number of fields: %d", len(fields))
	}

	result := Data{
		Country:      fields[0],
		BandWidth:    fields[1],
		ResponseTime: fields[2],
		Provider:     fields[3],
	}

	v4, err := strconv.ParseFloat(fields[4], 32)
	if err != nil {
		return Data{}, fmt.Errorf("invalid response time: %q", fields[4])
	}
	result.ConnectionStability = float32(v4)

	v5, err := strconv.Atoi(fields[5])
	if err != nil {
		return Data{}, fmt.Errorf("invalid ttfb: %v", fields[5])
	}
	result.Ttfb = v5

	v6, err := strconv.Atoi(fields[6])
	if err != nil {
		return Data{}, fmt.Errorf("invalid ttfb: %v", fields[6])
	}
	result.VoicePurity = v6

	v7, err := strconv.Atoi(fields[7])
	if err != nil {
		return Data{}, fmt.Errorf("invalid ttfb: %v", fields[7])
	}
	result.MedianOfCallsTime = v7

	if err := result.validate(); err != nil {
		return Data{}, err
	}

	return result, nil
}

func new(fileName string) (Set, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result Set
	var line string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()

		if d, err := decodeCSV(line); err == nil {
			if err := d.validate(); err != nil {
				return nil, err
			}
			result = append(result, d)
		}
	}

	return result, nil
}

func Result(fileName string) (Set, error) {
	result, err := new(fileName)
	if err != nil {
		return nil, err
	}

	return result, nil
}
