package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
)

func main() {
	ipLimiter := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	globalLimiter := NewConcurrentLimiter(10)

	http.Handle("/", globalLimiter.LimitConcurrentRequests(ipLimiter, HelloHandler))
	http.ListenAndServe(":8080", nil)
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

type ConcurrentLimiter struct {
	max     int
	current int
	mut     sync.Mutex
}

func NewConcurrentLimiter(limit int) *ConcurrentLimiter {
	return &ConcurrentLimiter{
		max: limit,
	}
}

func (limiter *ConcurrentLimiter) LimitConcurrentRequests(lmt *limiter.Limiter,
	handler func(http.ResponseWriter, *http.Request)) http.Handler {

	middle := func(w http.ResponseWriter, r *http.Request) {

		limiter.mut.Lock()
		maxHit := limiter.current == limiter.max

		if maxHit {
			limiter.mut.Unlock()
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		limiter.current += 1
		limiter.mut.Unlock()

		defer func() {
			limiter.mut.Lock()
			limiter.current -= 1
			limiter.mut.Unlock()
		}()

		// There's no rate-limit error, serve the next handler.
		handler(w, r)
	}

	return tollbooth.LimitHandler(lmt, http.HandlerFunc(middle))
}
