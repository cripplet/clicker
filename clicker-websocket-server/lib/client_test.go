package cc_websocket_server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"testing"
)

func TestStartClient(t *testing.T) {
	ts := NewMockServer(t, 8080, WebsocketHandler)

	ts.Start()
	defer ts.Close()

	c, r, e := websocket.DefaultDialer.Dial(strings.Replace(ts.URL, "http", "ws", 1), nil)
	defer c.Close()

	if e != nil {
		t.Fatalf("Unexpected error raised when attempting to establish websocket: %v", e)
	}
	if r.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("Unexpected HTTP status code returned: %d", r.StatusCode)
	}

	request := CommandRequest{}
	response := CommandResponse{}
	c.WriteJSON(&request)
	c.ReadJSON(&response)
	if response.Error.ErrorCode != ERROR_TYPE_INVALID_REQUEST {
		t.Errorf("Unexpected response status: %d != %d", response.Error.ErrorCode != ERROR_TYPE_INVALID_REQUEST)
	}
}
