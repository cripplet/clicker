package cc_rest_lib

import (
	"encoding/json"
	"fmt"
	"github.com/cripplet/clicker/db"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func NewMockServer(t *testing.T, p int, h func(http.ResponseWriter, *http.Request)) *httptest.Server {
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

	s := cc_fb.FBGameState{}
	json.Unmarshal(respRec.Body.Bytes(), &s)
	if s.ID == "" {
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

func TestNewGameHandlerClick(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	s := cc_fb.FBGameState{}
	json.Unmarshal(respRec.Body.Bytes(), &s)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/cookie/click", s.ID), nil)
	respRec = httptest.NewRecorder()
	http.HandlerFunc(ClickHandler).ServeHTTP(respRec, req)

	s, err := cc_fb.LoadGameState(s.ID)
	if err != nil {
		t.Errorf("Unexpected error when loading game: %v", err)
	}

	if s.GameData.NCookies == 0 {
		t.Error("Zero cookies when expecting more")
	}
}
