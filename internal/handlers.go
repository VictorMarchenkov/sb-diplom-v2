// Package internal  handlers to treat files and http queries
package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sb-diplom-v2/config"
	"sb-diplom-v2/internal/logger"
	"sb-diplom-v2/pkg"
	cfg2 "sb-diplom-v2/pkg/cfgPath"
	"sort"
	"strconv"
)

var (
	errFileOpen = fmt.Errorf("error reading data file")
	errFetchUrl = ""
)

// HandlerHTTP method of StatusResult for treating http queries.
func (t *StatusResult) HandlerHTTP(w http.ResponseWriter, r *http.Request) {
	t.Status = true

	MMS, err := GetMmsData(w, r)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Println("error when reading mms data: %v", err)
	}
	t.Data.MMS = MMS

	INCIDENT, err := GetIncidentData(w, r)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("error when reading incidents data: %v", err)
	}
	t.Data.Incident = INCIDENT

	SUPPORT, err := GetSupportServiceData(w, r)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("error when reading support service data: %v", err)
	}
	t.Data.Support = SUPPORT

	result, _ := json.Marshal(t)

	w.Write(result)
}

// HandlerFiles method of StatusResult for treating files.
func (t *StatusResult) HandlerFiles(cfg *cfg2.Root) {
	var err error

	t.Data.SMS, err = GetSmsData(cfg.CSV.Sms)
	if err != nil {
		logger.Logger.Errorf("%s", err)
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("%s", err)
	}
	t.Data.VoiceCall, err = GetVoiceService(cfg.CSV.Voice)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("%s", err)

	}
	t.Data.Email, err = GetEmailServiceData(cfg.CSV.Email)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("%s", err)

	}
	t.Data.Billing, err = GetBillingServiceData(cfg.CSV.Billing)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("%s", err)
	}
}

// GetMmsData collects mms data .
func GetMmsData(w http.ResponseWriter, r *http.Request) ([][]MMSData, error) {
	var (
		confT        Config
		tmpResult    []MMSData
		result       []MMSData
		sortedResult [][]MMSData
	)

	cfg, err := config.GetConfig()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return nil, err
	}
	json.Unmarshal(cfg, &confT)
	url := fmt.Sprintf("http://localhost:%d%s", confT.HTTP.ServicePort, confT.HTTP.Mms)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error parse %s: %v", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return [][]MMSData{}, nil
	}
	//	w.WriteHeader(200)
	rr, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error on response body for MMS service: %v", err)
		return [][]MMSData{}, nil
	}

	if err := json.Unmarshal(rr, &tmpResult); err != nil {
		fmt.Printf("error on decoding JSON response for MMS service: %s", err)
		return [][]MMSData{}, nil
	}
	for i := 0; i < len(tmpResult); i++ {
		if pkg.IsValidCountryCode(tmpResult[i].Country) && pkg.IsValidProvider(tmpResult[i].Provider) {
			result = append(result, tmpResult[i])
		} else {
			fmt.Println("something wrong... ", tmpResult[i].Country, ", or", tmpResult[i].Provider, " not valid")
		}
	}
	sortedResult = append(sortedResult, result)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Country < result[j].Country
	})

	resultCopy := append([]MMSData(nil), result...)
	sort.Slice(resultCopy, func(i, j int) bool {
		return resultCopy[i].Provider < resultCopy[j].Provider
	})
	sortedResult = append(sortedResult, resultCopy)

	for i := 0; i < len(sortedResult); i++ {
		for j := 0; j < len(sortedResult[i]); j++ {
			sortedResult[i][j].Country = pkg.CodeToName(sortedResult[i][j].Country)
		}
	}

	return sortedResult, nil
}

// GetSupportServiceData  collects support data.
func GetSupportServiceData(w http.ResponseWriter, r *http.Request) ([]int, error) {
	var (
		confT  Config
		result []SupportData
		report []int
	)

	cfg, err := config.GetConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	json.Unmarshal(cfg, &confT)
	url := fmt.Sprintf("http://localhost:%d%s", confT.HTTP.ServicePort, confT.HTTP.Support)

	res, err := http.Get(url)
	if err != nil {
		logger.Logger.Errorf("error parse url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return []int{}, nil
	}

	resu, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error on response body for support service: %v", err)
		return []int{}, nil
	}
	if err := json.Unmarshal(resu, &result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error on decoding JSON response for support service: %s", err)
		return []int{}, nil
	}

	allTickets := 0
	for i := 0; i < len(result); i++ {
		allTickets += result[i].ActiveTickets
	}

	timeChank := 60 / 18

	supportLoading := 1
	if allTickets <= 16.0 && allTickets > 8.0 {
		supportLoading = 2
	} else if allTickets > 16.0 {
		supportLoading = 3
	}
	report = []int{supportLoading, allTickets * timeChank}

	return report, nil
}

// GetIncidentData  collects incidents data.
func GetIncidentData(w http.ResponseWriter, r *http.Request) ([]IncidentData, error) {
	//return nil, fmt.Errorf("incident handler %s", "test error")
	var (
		confT  Config
		result []IncidentData
	)

	cfg, err := config.GetConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	json.Unmarshal(cfg, &confT)

	url := fmt.Sprintf("http://localhost:%d%s", confT.HTTP.ServicePort, confT.HTTP.Incident)

	res, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error parse url: %v", err)
		return []IncidentData{}, nil //[]IncidentData{}, nil
	}

	rr, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error on response body for insident service: %v", err)
		return []IncidentData{}, nil
	}

	if err := json.Unmarshal(rr, &result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error on decoding JSON response for insident service: %v", err)
		return []IncidentData{}, nil
	}
	//	w.WriteHeader(200)
	return result, nil
}

// GetSmsData  collects sms data.
func GetSmsData(path string) ([][]SMSData, error) {
	var err error
	var sms SMSData
	var result []SMSData
	var sortedResult [][]SMSData

	csv := pkg.ReadCSV(path)
	if csv == nil {
		err = errFileOpen
		return [][]SMSData{}, nil
	}

	for _, str := range csv {
		if len(str) == 4 {
			if pkg.IsValidCountryCode(str[0]) && pkg.IsValidProvider(str[3]) {
				sms.Country = str[0]
				sms.Bandwidth = str[1]
				sms.ResponseTime = str[2]
				sms.Provider = str[3]

				result = append(result, sms)
			}
		} else {
			//	fmt.Println("GetSmsData corrupted string", str)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Country < result[j].Country
	})
	resultCopy := append([]SMSData(nil), result...)
	sort.Slice(resultCopy, func(i, j int) bool {
		return resultCopy[i].Provider < resultCopy[j].Provider
	})
	sortedResult = append(sortedResult, result)
	sortedResult = append(sortedResult, resultCopy)
	for i := 0; i < len(sortedResult); i++ {
		for j := 0; j < len(sortedResult[i]); j++ {
			sortedResult[i][j].Country = pkg.CodeToName(sortedResult[i][j].Country)
		}
	}
	return sortedResult, err
}

// GetVoiceService  collects voice data.
func GetVoiceService(path string) ([]VoiceCallData, error) {
	var err error
	var voice VoiceCallData
	var result []VoiceCallData
	csv := pkg.ReadCSV(path)
	if csv == nil {
		err = errFileOpen
		result = append(result, VoiceCallData{}) //[]VoiceCallData{}
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
func GetEmailServiceData(path string) ([][]EmailData, error) {
	var err error
	var email EmailData
	var result []EmailData
	var resultCopy [][]EmailData

	csv := pkg.ReadCSV(path)
	if csv == nil {
		err = errFileOpen
		email = EmailData{"-", "-", 0}
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
				logger.Logger.Errorf("error reading sms: %s", err)
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

	return resultCopy, err
}

// GetBillingServiceData collects billing data  -.
func GetBillingServiceData(path string) (BillingData, error) {
	var err error
	var billing BillingData
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
		return BillingData{}, err
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
