// Package internal All entities of the project
package entities

import (
	"sb-diplom-v2/internal/entities/billing"
	"sb-diplom-v2/internal/entities/email"
	"sb-diplom-v2/internal/entities/incident"
	"sb-diplom-v2/internal/entities/mms"
	"sb-diplom-v2/internal/entities/sms"
	"sb-diplom-v2/internal/entities/support"
	"sb-diplom-v2/internal/entities/voice_call"
)

// SetStatusResult structure to collect al data.
type SetStatusResult struct {
	SMS       []sms.Set              `json:"sms"`
	MMS       []mms.Set              `json:"mms"`
	VoiceCall voice_call.Set         `json:"voice_call"`
	Email     map[string][]email.Set `json:"email"`
	Billing   billing.Data           `json:"billing"`
	Support   support.Set            `json:"support"`
	Incident  incident.Set           `json:"incident"`
}
