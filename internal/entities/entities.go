// Package internal All entities of the project
package entities

import (
	"sb-diplom-v2/internal/entities/email"
	"sb-diplom-v2/internal/entities/incident"
	"sb-diplom-v2/internal/entities/mms"
	"sb-diplom-v2/internal/entities/sms"
	"sb-diplom-v2/internal/entities/support"
)

// CSV

// SMSData describes the structure of data sms.
//type SMSData struct {
//	Country      string `json:"country"`
//	Bandwidth    string `json:"bandwidth"`
//	ResponseTime string `json:"response_time"`
//	Provider     string `json:"provider"`
//}

// VoiceCallData describes the data structure for call center services.
type VoiceCallData struct {
	Country             string  `json:"country"`
	Bandwidth           string  `json:"bandwidth"`
	ResponseTime        string  `json:"response_time"`
	Provider            string  `json:"provider"`
	ConnectionStability float32 `json:"connection_stability"`
	TTFB                int     `json:"ttfb"`
	VoicePurity         int     `json:"voice_purity"`
	MedianOfCallsTime   int     `json:"median_of_call_time"`
}

// EmailData describes the data structure for email services.
type EmailData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
}

// BillingData describes the data structure for billing services.
type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"recurring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}

// SetStatusResult structure to collect al data.
type SetStatusResult struct {
	SMS       []sms.Set              `json:"sms"`
	MMS       []mms.Set              `json:"mms"`
	VoiceCall []VoiceCallData        `json:"voice_call"`
	Email     map[string][]email.Set `json:"email"`
	Billing   BillingData            `json:"billing"`
	Support   support.Set            `json:"support"`
	Incident  incident.Set           `json:"incident"`
}
