package billing

import (
	"fmt"
	"os"
	"strconv"
)

type Data struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"recurring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}

func (d *Data) validate() error {
	return nil
}

func decodeCSV(str string) (Data, error) {

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

func Result(fileName string) (Data, error) {
	var result Data

	file, err := os.ReadFile(fileName)
	if err != nil {
		return result, err
	}
	result, err = decodeCSV(string(file))
	if err != nil {
		return result, err
	}

	return result, nil
}
