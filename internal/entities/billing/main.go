package billing

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Data struct {
	createCustomer bool
	purchase       bool
	payout         bool
	recurring      bool
	fraudControl   bool
	checkoutPage   bool
}

type Key int

const (
	createCustomer Key = 1 << iota
	purchase
	payout
	recurring
	fraudControl
	checkoutPage
)

func (d *Data) validate() error {

	if d.createCustomer != false && d.createCustomer != true {
		return fmt.Errorf("invalid value type for createCustomer: %T", d.createCustomer)
	}

	if d.purchase != false && d.purchase != true {
		return fmt.Errorf("invalid value type for purchase: %T", d.purchase)
	}

	if d.payout != false && d.payout != true {
		return fmt.Errorf("invalid value type for payout: %T", d.payout)
	}

	if d.fraudControl != false && d.fraudControl != true {
		return fmt.Errorf("invalid value type for fraudControl: %T", d.fraudControl)
	}

	if d.checkoutPage != false && d.checkoutPage != true {
		return fmt.Errorf("invalid value type for checkoutPage: %T", d.checkoutPage)
	}

	return nil
}

type Records []Data

func newFromString(str string) (Data, error) {

	result := Data{}

	if len(str) != 6 {
		return result, fmt.Errorf("error wrong length the string: %d", len(str))
	}

	d, err := strconv.ParseInt(str, 2, 32)
	if err != nil {
		return result, fmt.Errorf("error decoding of the string %v", err)
	}

	result.createCustomer = d&(1<<uint(0)) != 0
	result.purchase = d&(1<<uint(1)) != 0
	result.payout = d&(1<<uint(2)) != 0
	result.recurring = d&(1<<uint(3)) != 0
	result.fraudControl = d&(1<<uint(4)) != 0
	result.checkoutPage = d&(1<<uint(5)) != 0

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
