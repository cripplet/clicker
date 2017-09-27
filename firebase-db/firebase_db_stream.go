package firebase_db

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type EventType int

const (
	PUT EventType = iota
	PATCH
	KEEP_ALIVE
	CANCEL
	AUTH_REVOKED
)

var EVENT_TYPE_LOOKUP map[string]EventType = map[string]EventType{
	"put":          PUT,
	"patch":        PATCH,
	"keep-alive":   KEEP_ALIVE,
	"cancel":       CANCEL,
	"auth_revoked": AUTH_REVOKED,
}

// See https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#Event_stream_format.
type FirebaseDBStreamEvent struct {
	Event EventType
	Data  []byte
}

type FirebaseDBEventData struct {
	Path string `json:"path"`
	Data []byte `json:"data"`
}

func (self *FirebaseDBStreamEvent) GetEventData() (FirebaseDBEventData, error) {
	f := FirebaseDBEventData{}
	if self.Event == AUTH_REVOKED {
		return f, errors.New(fmt.Sprintf("Authorization revoked: %s", string(self.Data)))
	}
	err := json.Unmarshal(self.Data, &f)
	return f, err
}

func stream(
	c *http.Client,
	path string,
	query_parameters map[string]string) (*bufio.Scanner, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.Do(req)

	b_reader := bufio.NewScanner(resp.Body)

	return b_reader, resp.StatusCode, err
}
