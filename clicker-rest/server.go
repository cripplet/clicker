package main

import (
	"flag"
	"fmt"
	"github.com/cripplet/clicker/clicker-rest/lib"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"net/http"
	"time"
)

var rateLimiter *limiter.Limiter = tollbooth.NewLimiter(2, time.Second, nil)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Firebase DB environment")
	flag.Parse()

	http.Handle("/", tollbooth.LimitFuncHandler(rateLimiter, cc_rest_lib.GameRouter))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
