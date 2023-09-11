package sms

import (
	"bufio"
	"fmt"
	"os"
	"sb-diplom-v2/internal/entities/consts"
	"sort"
	"strconv"
	"strings"
)

type Data struct {
	Country      string `json:"country"`
	BandWidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

func (d *Data) validate() error {
	if len(d.Country) != 2 {
		return fmt.Errorf("invalid Country code length: %v", d.Country)
	}

	if d.BandWidth == "" {
		return fmt.Errorf("no bandwidth specified: %v", d.BandWidth)
	}

	_, err := strconv.ParseUint(d.BandWidth, 10, 8)
	if err != nil {
		return fmt.Errorf("invalid throughput: %q", d.BandWidth)
	}

	if d.ResponseTime == "" {
		return fmt.Errorf("no response time specified: %v", d.BandWidth)
	}

	_, err = strconv.ParseUint(d.ResponseTime, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid response time: %q", d.ResponseTime)
	}

	if _, ok := consts.ISOCountries[d.Country]; !ok {
		return fmt.Errorf("unknown Country code: %v", d.Country)
	}

	if d.Provider == "" {
		return fmt.Errorf("provider is empty")
	}

	if _, ok := consts.SMSProviders[d.Provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.Provider)
	}
	return nil
}

type Set []Data

func new(fileName string) (Set, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []Data
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

//
func newFromString(str string) (Data, error) {
	fields := strings.Split(str, ";")
	if len(fields) != 4 {
		return Data{}, fmt.Errorf("wrong number of fields: %d", len(fields))
	}
	result := Data{
		Country:      fields[0],
		BandWidth:    fields[1],
		ResponseTime: fields[2],
		Provider:     fields[3],
	}

	if err := result.validate(); err != nil {
		return Data{}, err
	}
	result.Country = consts.CountriesNames[result.Country]
	return result, nil
}

func SortedResult(fileName string) ([]Set, error) {
	s, err := new(fileName)
	if err != nil {
		return nil, err
	}

	result := make([]Set, 2)

	// sort by Country
	result[0] = s
	sort.Slice(s, func(i, j int) bool {
		return s[i].Country < s[j].Country
	})

	// sort by provider
	result[1] = make([]Data, len(s))
	copy(result[1], s)
	sort.Slice(result[1], func(i, j int) bool {
		return result[1][i].Provider < result[1][j].Provider
	})

	return result, nil
}
