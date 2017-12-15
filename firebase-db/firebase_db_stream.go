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

var EVENT_TYPE_LOOKUP map[EventType]string = map[EventType]string{
	PUT:          "put",
	PATCH:        "patch",
	KEEP_ALIVE:   "keep-alive",
	CANCEL:       "cancel",
	AUTH_REVOKED: "auth_revoked",
}

// EVENT_TYPE_REVERSE_LOOKUP is dynamically generated.
var EVENT_TYPE_REVERSE_LOOKUP map[string]EventType = map[string]EventType{}

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
	queryParameters map[string]string) (*bufio.Reader, int, error) {

	path += paramToURL(queryParameters)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.Do(req)

	r := bufio.NewReader(resp.Body)

	return r, resp.StatusCode, err
}

func doRead(r *bufio.Reader) {

}

func init() {
	for eventType, eventTypeString := range EVENT_TYPE_LOOKUP {
		EVENT_TYPE_REVERSE_LOOKUP[eventTypeString] = eventType
	}
}
