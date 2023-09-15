// Package internal  handlers to treat files and http queries
package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"sb-diplom-v2/internal/entities"
	"sb-diplom-v2/internal/entities/billing"
	"sb-diplom-v2/internal/entities/email"
	"sb-diplom-v2/internal/entities/incident"
	"sb-diplom-v2/internal/entities/mms"
	"sb-diplom-v2/internal/entities/sms"
	"sb-diplom-v2/internal/entities/support"
	"sb-diplom-v2/internal/entities/voice_call"
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

	t.Data.VoiceCall, err = voice_call.Result(t.cfg.CSV.Voice)
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

	t.Data.Billing, err = billing.Result(t.cfg.CSV.Billing)
	if err != nil {
		t.Status = false
		t.Error = fmt.Sprintf("%s", err)
		fmt.Printf("error Get Billing Service Data %s", err)
	}
}
