// Package internal All entities of the project
package internal

// CSV

// SMSData describes the structure of data sms.
type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

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

// HTTP

// MMSData describes the data structure for .
type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

// SupportData describes the data structure for .
type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

// IncidentData describes the data structure for incident service.
type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

// StatusResult structure to collect all information.
type StatusResult struct {
	Status bool            `json:"status"` // заполнен если все ОК, nil в противном случае
	Data   SetStatusResult `json:"data"`
	Error  string          `json:"error"` // пустая строка если все ОК, в противеом случае текст ошибки
}

// SetStatusResult structure to collect al data.
type SetStatusResult struct {
	SMS       [][]SMSData     `json:"sms"`
	MMS       [][]MMSData     `json:"mms"`
	VoiceCall []VoiceCallData `json:"voice_call"`
	Email     [][]EmailData   `json:"email"`
	Billing   BillingData     `json:"billing"`
	Support   []int           `json:"support"`
	Incident  []IncidentData  `json:"incident"`
}

// Config structure to represent config data.
type Config struct {
	CSV  Csv  `json:"csv"`
	HTTP Http `json:"httpreq"`
}

// Csv services.
type Csv struct {
	Sms     string `json:"sms"`
	Voice   string `json:"voice"`
	Email   string `json:"email"`
	Billing string `json:"billing"`
}

// Http services structure.
type Http struct {
	Mms         string `json:"mms"`
	Support     string `json:"support"`
	Incident    string `json:"incident"`
	ServerPort  int    `json:"server_port"`
	ServicePort int    `json:"service_port"`
}

type CacheRepository interface {
	Store(SetStatusResult) error
	Find(filter ValuesFilter) (SetStatusResult, error)
}

type ValuesFilter struct {
	CachedData SetStatusResult
	TimeStamp  string
}
