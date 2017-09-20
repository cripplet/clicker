package cc_websocket_server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SimpleType struct {
	Payload string
}

func SimpleWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	c, e := upgrader.Upgrade(w, r, nil)
	if e != nil {
		return
	}
	c.WriteJSON(SimpleType{Payload: "some-data"})
}

func NewMockServer(t *testing.T, p int, h func(http.ResponseWriter, *http.Request)) *httptest.Server {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(h))
	l, e := net.Listen("tcp", fmt.Sprintf("localhost:%d", p)) // http://blog.manugarri.com/how-to-mock-http-endpoints-in-golang/
	if e != nil {
		t.Fatalf("Unexpected error raised when constructing mock server: %v", e)
	}
	ts.Listener = l
	return ts
}

func TestSimpleWebsocketHandler(t *testing.T) {
	ts := NewMockServer(t, 8080, SimpleWebsocketHandler)

	ts.Start()
	defer ts.Close()

	c, r, e := websocket.DefaultDialer.Dial(strings.Replace(ts.URL, "http", "ws", 1), nil)
	if e != nil {
		t.Fatalf("Unexpected error raised when attempting to establish websocket: %v", e)
	}
	if r.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("Unexpected HTTP status code returned: %d", r.StatusCode)
	}

	actual := SimpleType{}
	c.ReadJSON(&actual)
	if actual.Payload != "some-data" {
		t.Errorf("Unexpected data received: %s", actual.Payload)
	}

}
