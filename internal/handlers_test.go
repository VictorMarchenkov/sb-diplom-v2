package internal

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"sb-diplom-v2/internal/logger"
	"testing"
)

func TestSmsHandler(t *testing.T) {

	var want [][]SMSData
	t.Run("is string with error deleted and result arranged", func(t *testing.T) {
		got, err := GetSmsData("./test-data/sms.data")
		if err != nil {
			logger.Logger.Warn("no fake data")
		}

		sms1 := SMSData{"Saint Barthélemy", "68", "1594", "Kildy"}
		sms2 := SMSData{"United States", "36", "1576", "Rond"}
		tmp1 := []SMSData{sms1, sms2}
		tmp2 := []SMSData{sms1, sms2}
		want = append(want, tmp1, tmp2)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
	t.Run("is some string deleted", func(t *testing.T) {
		got, err := GetSmsData("./test-data/sms.data")
		if err != nil {
			logger.Logger.Warn("no fake data")
		}

		sms1 := SMSData{"BL", "68", "1594", "Kildy"}
		sms2 := SMSData{"US", "36", "1576", "Rond"}
		tmp1 := []SMSData{sms1, sms2}
		tmp2 := []SMSData{sms1, sms2}
		want = append(want, tmp1, tmp2)

		if len(got) == len(want) {
			t.Errorf("len(got) '%d' len(want) '%d'", len(got), len(want))
		}
	})

}

func TestVoiceHandler(t *testing.T) {

	var want []VoiceCallData
	t.Run("is string with error deleted and result arranged", func(t *testing.T) {
		got, err := GetVoiceService("./test-data/voice.data")
		if err != nil {
			logger.Logger.Warn("no fake data")
		}
		voice1 := VoiceCallData{"BG", "40", "609", "E-Voice", 0.86, 160, 36, 5}
		voice2 := VoiceCallData{"DK", "11", "743", "JustPhone", 0.67, 82, 74, 41}

		want = append(want, voice1, voice2)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got '%v' want '%v'", got, want)
		}
	})
	t.Run("is some string deleted", func(t *testing.T) {
		got, err := GetVoiceService("./test-data/voice.data")
		if err != nil {
			logger.Logger.Warn("no fake data")
		}

		voice1 := VoiceCallData{"BG", "40", "609", "E-Voice", 0.86, 160, 36, 5}
		voice2 := VoiceCallData{"DK", "11", "743", "JustPhone", 0.67, 82, 74, 41}

		want = append(want, voice1, voice2)

		if len(got) == len(want) {
			t.Errorf("len(got) '%d' len(want) '%d'", len(got), len(want))
		}
	})

}

func TestEmailHandler(t *testing.T) {

	var want [][]EmailData
	t.Run("is error string deleted", func(t *testing.T) {
		got, err := GetEmailServiceData("./test-data/email.data")
		if err != nil {
			logger.Logger.Warn("no fake data")
		}
		sms1 := EmailData{"AT", "Hotmail", 487}
		tmp1 := []EmailData{sms1}
		want = append(want, tmp1)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got '%v' want '%v'", got, want)
		}
	})

}

func TestBillingHandler(t *testing.T) {

	t.Run("is string decoded correctly", func(t *testing.T) {
		got, err := GetBillingServiceData("./test-data/billing.data")
		if err != nil {
			logger.Logger.Warn("no fake data")
		}
		fmt.Printf("got %T %v\n", got, got)

		want := BillingData{true, true, false, false, true, false}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("got '%v' want '%v'", got, want)
		}
	})

}
