package cfg

import "fmt"

type Root struct {
	HTTPServer httpServer `json:"http_server"`
	CSV        csv        `json:"csv"`
}

func (r *Root) initiate() error {
	if err := r.HTTPServer.initiate(); err != nil {
		return fmt.Errorf("http_server: %w", err)
	}
	if err := r.CSV.initiate(); err != nil {
		return fmt.Errorf("csv: %w", err)
	}

	return nil
}
