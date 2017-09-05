package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	ccfb_config "github.com/cripplet/clicker/db/firebase/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"time"
)

type SomeData struct {
	ID       string
	Username string
	CPS      float64
}

func DoRequest(c *http.Client, s SomeData) {
	s_json, _ := json.Marshal(s)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/game/%s.json", ccfb_config.DB_CONFIG.BaseURL, s.ID),
		bytes.NewReader(s_json),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while HTTP req: %v", err))
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(resp.StatusCode, resp.ContentLength, fmt.Sprintf("%s", body))
	}

}

// See https://github.com/golang/oauth2/blob/d89af98d7c6bba047c5a2622f36bc14b8766df85/google/jwt_test.go#L22
func main() {
	var credentials string
	flag.StringVar(&credentials, "credentials", "", "JSON credentials file")
	flag.Parse()

	json_key, err := ioutil.ReadFile(credentials)
	if err != nil {
		panic(fmt.Sprintf("Attempted to open credentials file \"%s\", returned error: \"%v\"", credentials, err))
	}

	// See https://godoc.org/golang.org/x/oauth2/google#JWTConfigFromJSON
	config, _ := google.JWTConfigFromJSON(
		json_key,
		"https://www.googleapis.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email",
	)

	c := config.Client(oauth2.NoContext) // HTTP client with wrapped auth token headers.

	s := SomeData{
		ID:       "some-id",
		Username: "cripplet",
		CPS:      0,
	}

	DoRequest(c, s)

	time.Sleep(time.Second * 3700)

	DoRequest(c, s)
}
