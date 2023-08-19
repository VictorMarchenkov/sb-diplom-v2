package mms

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"sb-diplom-v2/internal/entities/consts"
)

type Data struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type Set []Data

func (d *Data) validate() error {
	if len(d.Country) != 2 {
		return fmt.Errorf("invalid country code length: %v", d.Country)
	}

	if _, ok := consts.ISOCountries[d.Country]; !ok {
		return fmt.Errorf("unknown country code: %v", d.Country)
	}

	if d.Provider == "" {
		return fmt.Errorf("provider is empty")
	}

	if _, ok := consts.MMSProviders[d.Provider]; !ok {
		return fmt.Errorf("unknown provider: %v", d.Provider)
	}
	return nil
}

func new(url string) (Set, error) {
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

func SortedResult(url string) ([]Set, error) {
	s, err := new(url)
	if err != nil {
		return nil, err
	}

	result := make([]Set, 2)

	// sort by country
	result[0] = s
	sort.Slice(s, func(i, j int) bool {
		return s[i].Country < s[j].Country
	})

	// sort by provider
	result[0] = s
	sort.Slice(s, func(i, j int) bool {
		return s[i].Country < s[j].Country
	})

	// sort by provider
	result[1] = make([]Data, len(s))
	copy(result[1], s)
	sort.Slice(result[1], func(i, j int) bool {
		return result[1][i].Provider < result[1][j].Provider
	})

	return result, nil
}
