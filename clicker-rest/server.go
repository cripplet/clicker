package main

import (
	"flag"
	"fmt"
	"github.com/cripplet/clicker/clicker-rest/lib"
	"github.com/cripplet/clicker/db/config"
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

	cc_fb_config.SetCCFirebaseConfig()
	if cc_fb_config.CC_FIREBASE_CONFIG.Environment != cc_fb_config.DEV {
		panic(fmt.Sprintf("Firebase environment is not %s", cc_fb_config.ENVIRONMENT_TYPE_LOOKUP[cc_fb_config.DEV]))
	}

	http.Handle("/", tollbooth.LimitFuncHandler(rateLimiter, cc_rest_lib.GameRouter))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
