package cc_websocket_server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
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

	r_bytes, _ := json.Marshal(&response)
	fmt.Printf("%s", string(r_bytes))
}
