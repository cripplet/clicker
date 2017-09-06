package cc_firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	cc_firebase_config "github.com/cripplet/clicker/db/firebase/config"
	"io/ioutil"
	"net/http"
)

type Session struct { // TODO(cripplet): Add FirebaseStruct.
	ID string
}

type User struct {
	ID        string
	UserToken string
}

func CloneSession(c *http.Client, e Environment, id string) (*Session, error) { // TODO(cripplet): Define GenerateRandomID with stateful cache.
	s, err := ReadSession(c, e, id)
	if err != nil {
		return nil, err
	}

	s.ID += "-clone"
	return CreateSession(c, e, s.ID)
}

func CreateSession(c *http.Client, e Environment, id string) (*Session, error) { // TODO(cripplet): Use GenerateRandomID.
	s := Session{
		ID: id,
	}

	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf(
			"%s/%s/game/%s.json",
			cc_firebase_config.DB_CONFIG.BaseURL,
			ENVIRONMENT_LOOKUP[e],
			id,
		),
		bytes.NewReader(j),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &s)
	return &s, nil
}

func ReadSession(c *http.Client, e Environment, id string) (*Session, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/game/%s.json",
			cc_firebase_config.DB_CONFIG.BaseURL,
			ENVIRONMENT_LOOKUP[e],
			id,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	s := Session{}
	json.Unmarshal(body, &s)
	return &s, nil
}

func NewSession(id string) *Session {
	if id != "" {
	}
	return nil
}
