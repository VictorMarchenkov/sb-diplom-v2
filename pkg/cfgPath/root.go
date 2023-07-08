package cfgPath

import "fmt"

type Root struct {
	HTTPServer  httpServer  `json:"http_server"`
	CSV         csv         `json:"csv"`
	HTTPService httpService `json:"http_service"`
}

func (r *Root) initiate() error {

	if err := r.HTTPServer.initiate(); err != nil {
		return fmt.Errorf("http_server: %w", err)
	}

	if err := r.CSV.initiate(); err != nil {
		return fmt.Errorf("csv: %w", err)
	}

	if err := r.HTTPService.initiate(); err != nil {
		return fmt.Errorf("http service: %w", err)
	}

	return nil
}
