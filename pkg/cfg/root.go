package cfg

import "fmt"

type Root struct {
	HTTPServer httpServer `json:"http_server"`
}

func (r *Root) initiate() error {
	if err := r.HTTPServer.initiate(); err != nil {
		return fmt.Errorf("http_server: %w", err)
	}

	return nil
}
