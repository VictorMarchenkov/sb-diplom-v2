package pkg

import (
	"fmt"
	"net/http"
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("PanicMiddleware recovered: %s", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
