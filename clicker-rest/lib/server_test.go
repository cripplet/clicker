package cc_rest_lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cripplet/clicker/db"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func getGameIDFromPath(p string) string {
	r := regexp.MustCompile(`/(?P<gameID>[\w-]*)\.json$`)
	return regexpMatchNamedGroups(r, p)["gameID"]
}

func newMockServer(t *testing.T, p int, h func(http.ResponseWriter, *http.Request)) *httptest.Server {
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(h))
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", p)) // http://blog.manugarri.com/how-to-mock-http-endpoints-in-golang/
	if err != nil {
		t.Fatalf("Unexpected error raised when constructing mock server: %v", err)
	}
	testServer.Listener = listener
	return testServer
}

func TestRegexpMatchNamedGroups(t *testing.T) {
	r := regexp.MustCompile(`(?P<first>a)(?P<second>b)`)
	m := regexpMatchNamedGroups(r, "ab")
	if len(m) != 2 {
		t.Errorf("Length of map mismatch: %d != %d", len(m), 2)
	}

	if m["first"] != "a" {
		t.Errorf("Incorrect match: %s != %s", m["first"], "a")
	}
	if m["second"] != "b" {
		t.Errorf("Incorrect match: %s != %s", m["second"], "b")
	}
}

func TestNewGameHandler(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP error code %d: %s", respRec.Result().StatusCode, string(respRec.Body.Bytes()))
	}

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)
	if g.Path == "" {
		t.Error("Game ID was not set when creating new game")
	}
}

func TestNewGameHandlerInvalidMethod(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodGet, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected HTTP error code %d: %s", respRec.Result().StatusCode, string(respRec.Body.Bytes()))
	}
}

func TestClickHandler(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	clickRequest := ClickRequest{
		NTimes: 10,
	}
	clickRequestJSON, _ := json.Marshal(&clickRequest)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/cookie/click", getGameIDFromPath(g.Path)), bytes.NewReader(clickRequestJSON))
	respRec = httptest.NewRecorder()
	http.HandlerFunc(ClickHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusNoContent {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusNoContent)
	}

	s, _, err := cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	if err != nil {
		t.Errorf("Unexpected error when loading game: %v", err)
	}

	if s.GameData.NCookies == 0 {
		t.Error("Zero cookies when expecting more")
	}
}
