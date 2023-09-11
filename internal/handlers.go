// Package internal  handlers to treat files and http queries
package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sb-diplom-v2/internal/entities"
	"sb-diplom-v2/internal/entities/email"
	"sb-diplom-v2/internal/entities/incident"
	"sb-diplom-v2/internal/entities/mms"
	"sb-diplom-v2/internal/entities/sms"
	"sb-diplom-v2/internal/entities/support"
	"sb-diplom-v2/pkg"
	"sb-diplom-v2/pkg/configs"
	"sb-diplom-v2/pkg/logger"
)

var (
	errFileOpen = errors.New("error reading data file")
	errFetchUrl = errors.New("error fetching url")
)

// StatusResult structure to collect all information.
type StatusResult struct {
	Status bool                     `json:"status"` // заполнен если все ОК, nil в противном случае
	Data   entities.SetStatusResult `json:"data"`
	Error  string                   `json:"error"` // пустая строка если все ОК, в противеом случае текст ошибки
	l      logger.Logger
	cfg    *configs.Root
}

func NewStatusResult(cfg *configs.Root) *StatusResult {
	return &StatusResult{l: logger.New("status-result"), cfg: cfg}
}

// HandlerHTTP method of StatusResult for treating http queries.
func (t *StatusResult) HandlerHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	t.Data.MMS, err = mms.SortedResult(t.cfg.HTTPService.MMSURL)
	if err != nil {
		t.Status = false
		t.Error = err.Error()
		http.Error(w, "error read of mms data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t.Data.Incident, err = incident.Result(t.cfg.HTTPService.IncidentURL)
	if err != nil {
		t.Status = false
		t.Error = err.Error()
		http.Error(w, "error read of incident data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t.Data.Support, err = support.Result(t.cfg.HTTPService.SupportURL)
	if err != nil {
		t.Status = false
		t.Error = err.Error()
		http.Error(w, "error read of support data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, _ := json.Marshal(t)

	w.Write(result)
}

// HandlerFiles method of StatusResult for treating files.
func (t *StatusResult) HandlerFiles(cfg *configs.Root) {
	var err error

	t.Data.SMS, err = sms.SortedResult(t.cfg.CSV.Sms)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("error Get Billing Service Data %s", err)
	}

	t.Data.VoiceCall, err = GetVoiceService(t.cfg.CSV.Voice)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("error Get Voice Service Data %s", err)

	}

	t.Data.Email, err = email.SortedResult(t.cfg.CSV.Email)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("error Get Email Service Data %s", err)

	}

	t.Data.Billing, err = GetBillingServiceData(cfg.CSV.Billing)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("error Get Billing Service Data %s", err)
	}
}

// GetVoiceService  collects voice data.
func GetVoiceService(path string) ([]entities.VoiceCallData, error) {
	var err error
	var voice entities.VoiceCallData
	var result []entities.VoiceCallData
	csv := pkg.ReadCSV(path)
	if csv == nil {
		err = errFileOpen
		result = append(result, entities.VoiceCallData{}) //[]VoiceCallData{}
		return result, nil
	}
	for _, str := range csv {

		if len(str) == 8 {
			if pkg.IsValidCountryCode(str[0]) && pkg.IsValidVoiceProvider(str[3]) {
				cstability, err0 := strconv.ParseFloat(str[4], 32)
				ttfb, err1 := strconv.Atoi(str[5])
				vpurity, err2 := strconv.Atoi(str[6])
				mtime, err3 := strconv.Atoi(str[7])
				if err0 == nil && err1 == nil && err2 == nil && err3 == nil {
					voice.Country = str[0]
					voice.Bandwidth = str[1]
					voice.ResponseTime = str[2]
					voice.Provider = str[3]
					voice.ConnectionStability = float32(cstability)
					voice.TTFB = ttfb
					voice.VoicePurity = vpurity
					voice.MedianOfCallsTime = mtime

					result = append(result, voice)
				}
			}
		}
	}

	return result, err
}

// GetEmailServiceData  collects email data.
func GetEmailServiceData(path string) ([][]entities.EmailData, error) {
	var err error
	var email entities.EmailData
	var result []entities.EmailData
	var resultCopy [][]entities.EmailData

	csv := pkg.ReadCSV(path)
	if csv == nil {
		err = errFileOpen
		email = entities.EmailData{"-", "-", 0}
		result = append(result, email)
		result = append(result, email)
		result = append(result, email)
		resultCopy = append(resultCopy, result)
		return resultCopy, nil
	}
	for _, str := range csv {
		if len(str) == 3 {
			dtime, err := strconv.Atoi(str[2])
			if err != nil {
				fmt.Printf("error reading sms: %s", err)
			}
			if pkg.IsValidCountryCode(str[0]) && pkg.IsValidEmailProvider(str[1]) && err == nil {
				email.Country = str[0]
				email.Provider = str[1]
				email.DeliveryTime = dtime
				result = append(result, email)
			}
		}
	}

	if len(result) == 1 {
		resultCopy = append(resultCopy, result[0:1])
	} else {
		resultCopy = append(resultCopy, result[0:3])
		if len(result) > 3 {
			resultCopy = append(resultCopy, result[len(result)-4:len(result)-1])
		}
	}
	fmt.Println("\nresultCopy", resultCopy)
	return resultCopy, err
}

// GetBillingServiceData collects billing data  -.
func GetBillingServiceData(path string) (entities.BillingData, error) {
	var err error
	var billing entities.BillingData
	type Key int
	const (
		CreateCustomer Key = 1 << iota
		Purchase
		Payout
		Recurring
		FraudControl
		CheckoutPage
	)

	csv := pkg.ReadCSV(path)
	if csv == nil {
		err = errFileOpen
		return entities.BillingData{}, err
	}
	for _, str := range csv {
		d, err := strconv.ParseInt(str[0], 2, 32)
		if err == nil && len(str[0]) == 6 {
			billing.CreateCustomer = d&(1<<uint(0)) != 0
			billing.Purchase = d&(1<<uint(1)) != 0
			billing.Payout = d&(1<<uint(2)) != 0
			billing.Recurring = d&(1<<uint(3)) != 0
			billing.FraudControl = d&(1<<uint(4)) != 0
			billing.CheckoutPage = d&(1<<uint(5)) != 0
		}
	}
	return billing, nil
}
