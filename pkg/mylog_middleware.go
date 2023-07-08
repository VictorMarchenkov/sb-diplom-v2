package pkg

import (
	"net/http"
	"sb-diplom-v2/pkg/logger"
	"time"
)

// AccessLogMiddleware middleware handlers wrapper
func AccessLogMiddleware(next http.Handler) http.Handler {
	l := logger.New("access-log")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		l.Info("[%s] call from address: %s, requested url: %s, time to respone: %s\n",
			r.Method, r.Host, r.URL.Path, time.Since(start))
	})
}
