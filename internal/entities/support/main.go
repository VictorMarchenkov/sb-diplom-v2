package support

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Data struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

type Set []int

func (d *Data) validate() error {
	if d.Topic == "" {
		return fmt.Errorf("topic is empty")
	}

	if d.ActiveTickets < 0 {
		return fmt.Errorf("error reading active tickets")
	}

	return nil
}

func new(url string) ([]Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: set timeout from config
	defer cancel()

	c := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var result []Data
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	// validation
	for _, r := range result {
		if err := r.validate(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func Result(url string) ([]int, error) {
	result, err := new(url)
	if err != nil {
		return []int{0, 0}, nil
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
	report := []int{supportLoading, allTickets * timeChank}

	return report, nil
}
