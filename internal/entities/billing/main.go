package billing

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Data struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

func (d *Data) validate() error {
	return nil
}

type Records []Data

func newFromString(str string) (Data, error) {

	result := Data{}

	if len(str) != 6 {
		return result, fmt.Errorf("wrong length the string: %d", len(str))
	}

	d, err := strconv.ParseInt(str, 2, 32)
	if err != nil {
		return result, fmt.Errorf("decoding of the string %v", err)
	}

	result.CreateCustomer = d&(1<<uint(0)) != 0
	result.Purchase = d&(1<<uint(1)) != 0
	result.Payout = d&(1<<uint(2)) != 0
	result.Recurring = d&(1<<uint(3)) != 0
	result.FraudControl = d&(1<<uint(4)) != 0
	result.CheckoutPage = d&(1<<uint(5)) != 0

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
