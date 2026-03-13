package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	requests int
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func getIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func RateLimiter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := getIP(r)

		mu.Lock()

		v, exists := visitors[ip]

		if !exists {
			visitors[ip] = &visitor{requests: 1, lastSeen: time.Now()}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSeen) > time.Minute {
			v.requests = 0
			v.lastSeen = time.Now()
		}

		v.requests++

		if v.requests > 100 {
			mu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
