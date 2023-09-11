package email

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"sb-diplom-v2/internal/entities/consts"
)

type Data struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
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

	if _, ok := consts.EmailProviders[d.Provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.Provider)
	}

	var number int

	if reflect.ValueOf(d.DeliveryTime).Type() != reflect.TypeOf(number) {
		return errors.New("wrong type for deliveryTime")
	}

	return nil
}

type Set []Data

func newFromString(str string) (Data, error) {
	fields := strings.Split(str, ";")
	if len(fields) != 3 {
		return Data{}, fmt.Errorf("wrong number of fields: %d", len(fields))
	}

	var v int

	if fields[2] != "" {
		v, _ = strconv.Atoi(fields[2])

	}
	result := Data{
		Country:      fields[0],
		Provider:     fields[1],
		DeliveryTime: v,
	}

	if err := result.validate(); err != nil {
		return result, err
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
		if d, err := newFromString(line); err == nil {
			result = append(result, d)
		}
	}

	return result, nil
}

type ResultData map[string][]Set

func SortedResult(fileName string) (map[string][]Set, error) {
	s, err := new(fileName)
	if err != nil {
		return nil, err
	}

	preFinal := make(map[string]Set)

	for _, c := range s {
		preFinal[c.Country] = append(preFinal[c.Country], c)
	}

	final := make(ResultData)

	for country := range preFinal {
		sort.Slice(preFinal[country], func(i, j int) bool {
			return int(s[i].DeliveryTime) < int(s[j].DeliveryTime)
		})
		/*
			if len(result) == 1 {
				resultCopy = append(resultCopy, result[0:1])
			} else {
				resultCopy = append(resultCopy, result[0:3])
				if len(result) > 3 {
					resultCopy = append(resultCopy, result[len(result)-4:len(result)-1])
				}
			}
		*/
		final[country] = append(final[country], preFinal[country][0:3])
		final[country] = append(final[country], preFinal[country][len(preFinal[country])-4:len(preFinal[country])-1])
	}

	return final, nil
}
