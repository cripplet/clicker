package firebase_db

import (
	"bytes"
	"encoding/json"
	// "fmt"
	// "net/http"
	"testing"
)

func TestGetEventData(t *testing.T) {
	event_data := FirebaseDBEventData{
		Path: "/some/path",
		Data: []byte("0"),
	}
	expected, _ := json.Marshal(event_data)

	stream_event := FirebaseDBStreamEvent{
		Event: PUT,
		Data:  expected,
	}

	actual_event_data, err := stream_event.GetEventData()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actual_event_data.Path != "/some/path" {
		t.Errorf("Unexpected path: %s != %s", actual_event_data.Path, "/some/path")
	}

	if !bytes.Equal(actual_event_data.Data, []byte("0")) {
		t.Errorf("Unexpected data: %s != %s", string(actual_event_data.Data), string([]byte("0")))
	}
}

func TestGetEventDataAuthRevoked(t *testing.T) {
	stream_event := FirebaseDBStreamEvent{
		Event: AUTH_REVOKED,
		Data:  []byte("Error message"),
	}
	_, e := stream_event.GetEventData()
	if e == nil {
		t.Errorf("Expected error was not raised")
	}
}

/*
func TestStream(t *testing.T) {
	ResetEnvironment(t)

	c := &http.Client{}

	sc, status_code, _ := stream(
		c,
		fmt.Sprintf("%s/test_stream.json", project_root),
		map[string]string{},
	)

	if status_code != http.StatusOK {
		t.Errorf("HTTP Error: %d", status_code)
	}

	var e FirebaseDBStreamEvent = FirebaseDBStreamEvent{}

	for sc.Scan() {
		line := sc.Bytes()
		fmt.Println("%v", line)
		if len(line) == 0 {
			fmt.Printf("%s", e.Event)
			fmt.Printf("%s", string(e.Data))
			e = FirebaseDBStreamEvent{}
			continue
		}
		line_tokens := bytes.SplitN(line, []byte(": "), 2)
		switch string(line_tokens[0]) {
			case "event":
				e.Event = EVENT_TYPE_REVERSE_LOOKUP[string(line_tokens[1])]
				break
			case "data":
				e.Data = append(e.Data, line_tokens[1]...)
				break
			default:
				fmt.Printf("No match: %s", string(line))
		}
	}
}
*/
