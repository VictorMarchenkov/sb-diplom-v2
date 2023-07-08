package pkg

import (
	"fmt"
	"net/http"
	"time"
)

// AccessLogMiddleware middleware handlers wrapper
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] call from address: %s, requested url: %s, time to respone: %s\n",
			r.Method, r.Host, r.URL.Path, time.Since(start))
	})
}
