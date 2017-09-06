package cc_firebase

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

func NewFirebaseClient(credentials string) (*http.Client, error) {
	data, err := ioutil.ReadFile(credentials)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(
		data,
		"https://www.googleapis.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email",
	)
	if err != nil {
		return nil, err
	}

	client := conf.Client(oauth2.NoContext)
	return client, nil
}
