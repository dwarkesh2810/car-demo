package middleware

import (
	"car_demo/conf"
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

var (
	accessCount = make(map[string]int)
	mutex       = &sync.Mutex{}
	blocked     = make(map[string]bool)
	unblocked   = make(map[string]int64)
)

func RateLimiter(ctx *context.Context) {
	// Get IP address of the client
	ip := ctx.Input.IP()

	// Limit requests from an IP address
	mutex.Lock()
	defer mutex.Unlock()

	accessCount[ip]++
	if accessCount[ip] > conf.EnvConfig.RateLimiter {
		blocked[ip] = true
		if blocked[ip] && unblocked[ip] > 0 {
			if unblocked[ip] < time.Now().Unix() {
				accessCount[ip] = 0
				blocked[ip] = false
				unblocked[ip] = 0
				return
			}
		}
		ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
		ctx.ResponseWriter.Write([]byte("Rate limit exceeded. Please try again later."))

		if blocked[ip] {
			unblocked[ip] = int64(time.Now().Add(60 * time.Second).Unix())
		}
		return
	}
}
