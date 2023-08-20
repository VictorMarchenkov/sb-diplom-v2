package sms

import (
	"fmt"
	"sb-diplom-v2/internal/entities/consts"
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

	if _, ok := consts.ISOCountries[d.country]; ok {
		return fmt.Errorf("unknown country code: %v", d.country)
	}

	if d.provider == "" {
		return fmt.Errorf("provider is empty")
	}

	if _, ok := consts.SMSProviders[d.provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.provider)
	}
	return nil
}

type Records []Data

//
//func newFromString(str string) (Data, error) {
//	fields := strings.Split(str, ";")
//	if len(fields) != 4 {
//		return Data{}, fmt.Errorf("wrong number of fields: %d", len(fields))
//	}
//	result := Data{
//		country:  fields[0],
//		provider: fields[3],
//	}
//
//	v, err := strconv.ParseUint(fields[1], 10, 8)
//	if err != nil {
//		return Data{}, fmt.Errorf("invalid throughput: %q", fields[1])
//	}
//	result.bandWidth = uint8(v)
//
//	v, err = strconv.ParseUint(fields[2], 10, 64)
//	if err != nil {
//		return Data{}, fmt.Errorf("invalid response time: %q", fields[2])
//	}
//	result.responseTime = uint(v)
//
//	if err := result.validate(); err != nil {
//		return Data{}, err
//	}
//
//	return result, nil
//}

//func NewFromFile(fileName string) (Records, error) {
//	file, err := os.Open(fileName)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	var result Records
//	var line string
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		line = scanner.Text()
//		if d, err := newFromString(line); err == nil {
//			result = append(result, d)
//		}
//	}
//
//	return result, nil
//}
