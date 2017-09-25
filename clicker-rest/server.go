package main

import (
	"github.com/cripplet/clicker/clicker-rest/lib"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"net/http"
	"time"
)

var rateLimiter *limiter.Limiter = tollbooth.NewLimiter(2, time.Second, nil)

func main() {
	http.Handle("/", tollbooth.LimitFuncHandler(rateLimiter, cc_rest_lib.GameRouter))
	http.ListenAndServe(":8080", nil)
}
