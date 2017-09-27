package cc_rest_lib

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cripplet/clicker/db"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/lib"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
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

func TestGameRouterNonexistentPath(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/nonexistent-path", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(GameRouter).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Unexpected HTTP error code: %d != %d", respRec.Result().StatusCode, http.StatusNotFound)
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

	nCookies := int(1e5)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	clickRequest := ClickRequest{
		Events: []time.Time{},
	}

	s, _, err := cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	h := s.Metadata.ClickHash
	for i := 0; i < nCookies; i++ {
		newHashBytes := sha256.Sum256(h)
		h = newHashBytes[:]
		clickRequest.Events = append(clickRequest.Events, time.Now())
	}
	clickRequest.Hash = h

	clickTime := s.Metadata.ClickTime
	clickHash := h

	clickRequestJSON, _ := json.Marshal(&clickRequest)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/cookie/click/", getGameIDFromPath(g.Path)), bytes.NewReader(clickRequestJSON))
	respRec = httptest.NewRecorder()
	http.HandlerFunc(ClickHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusNoContent {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusNoContent)
	}

	s, _, err = cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	if err != nil {
		t.Errorf("Unexpected error when loading game: %v", err)
	}

	if s.GameData.NCookies == 0 {
		t.Error("Zero cookies when expecting more")
	}

	if !clickTime.Before(s.Metadata.ClickTime) {
		t.Error("Click time was not set upon a successful click request")
	}

	if !bytes.Equal(clickHash, s.Metadata.ClickHash) {
		t.Error("Click hash was not set upon a successful click request")
	}
}

func TestClickHandlerInvalidTime(t *testing.T) {
	n := time.Now()
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	s, _, _ := cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	h := s.Metadata.ClickHash
	newHashBytes := sha256.Sum256(h)
	h = newHashBytes[:]

	clickRequest := ClickRequest{
		Events: []time.Time{
			n,
		},
		Hash: h,
	}

	clickRequestJSON, _ := json.Marshal(&clickRequest)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/cookie/click/", getGameIDFromPath(g.Path)), bytes.NewReader(clickRequestJSON))
	respRec = httptest.NewRecorder()
	http.HandlerFunc(ClickHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestClickHandlerInvalidHash(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	clickRequest := ClickRequest{
		Events: []time.Time{
			time.Now(),
		},
		Hash: []byte("invalid-hash"),
	}
	clickRequestJSON, _ := json.Marshal(&clickRequest)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/cookie/click/", getGameIDFromPath(g.Path)), bytes.NewReader(clickRequestJSON))
	respRec = httptest.NewRecorder()
	http.HandlerFunc(ClickHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusBadRequest)
	}

	s, _, _ := cc_fb.LoadGameState(getGameIDFromPath(g.Path))

	if s.GameData.NCookies != 0 {
		t.Error("Cookies found in game where none should be")
	}
}

func TestMineHandler(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	s, eTag, _ := cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	s.GameData.NBuildings[cookie_clicker.BUILDING_TYPE_CURSOR] = 1
	cc_fb.SaveGameState(s, eTag)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/cookie/mine/", getGameIDFromPath(g.Path)), nil)
	respRec = httptest.NewRecorder()
	http.HandlerFunc(MineHandler).ServeHTTP(respRec, req)

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

func TestBuildingHandlerNonexistentBuilding(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/some-id/building/nonexistent-building-type", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(BuildingHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusNotFound)
	}
}

func TestBuildingHandlerInsufficientFunds(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/building/%s/", getGameIDFromPath(g.Path), cookie_clicker.BUILDING_TYPE_LOOKUP[cookie_clicker.BUILDING_TYPE_CURSOR]), nil)
	respRec = httptest.NewRecorder()
	http.HandlerFunc(BuildingHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusPaymentRequired {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusPaymentRequired)
	}
}

func TestBuildingHandler(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	s, eTag, _ := cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	game := cookie_clicker.NewGameState()
	game.Load(s.GameData)
	var i float64
	for i = 0; i < game.GetBuildings()[cookie_clicker.BUILDING_TYPE_CURSOR].GetCost(1); i++ {
		game.Click()
	}
	s.GameData = game.Dump()
	cc_fb.SaveGameState(s, eTag)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/building/%s/", getGameIDFromPath(g.Path), cookie_clicker.BUILDING_TYPE_LOOKUP[cookie_clicker.BUILDING_TYPE_CURSOR]), nil)
	respRec = httptest.NewRecorder()
	http.HandlerFunc(BuildingHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusNoContent {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusNoContent)
	}

	s, _, _ = cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	if s.GameData.NBuildings[cookie_clicker.BUILDING_TYPE_CURSOR] != 1 {
		t.Errorf("Game state does not reflect building bought: %d buildings found", s.GameData.NBuildings[cookie_clicker.BUILDING_TYPE_CURSOR])
	}
}

func TestUpgradeHandlerNonexistentUpgrade(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/some-id/upgrade/nonexistent-upgrade", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(UpgradeHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusNotFound {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusNotFound)
	}
}

func TestUpgradeHandlerInsufficientHandler(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/upgrade/%s/", getGameIDFromPath(g.Path), cookie_clicker.UPGRADE_ID_LOOKUP[cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER]), nil)
	respRec = httptest.NewRecorder()
	http.HandlerFunc(UpgradeHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusPaymentRequired {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusPaymentRequired)
	}
}

func TestUpgradeHandler(t *testing.T) {
	cc_fb.ResetEnvironment(t)

	req, _ := http.NewRequest(http.MethodPost, "/game/", nil)
	respRec := httptest.NewRecorder()
	http.HandlerFunc(NewGameHandler).ServeHTTP(respRec, req)

	g := NewGameResponse{}
	json.Unmarshal(respRec.Body.Bytes(), &g)

	s, eTag, _ := cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	game := cookie_clicker.NewGameState()
	game.Load(s.GameData)
	var i float64
	for i = 0; i < game.GetUpgrades()[cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER].GetCost(game); i++ {
		game.Click()
	}
	s.GameData = game.Dump()
	cc_fb.SaveGameState(s, eTag)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("/game/%s/upgrade/%s/", getGameIDFromPath(g.Path), cookie_clicker.UPGRADE_ID_LOOKUP[cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER]), nil)
	respRec = httptest.NewRecorder()
	http.HandlerFunc(UpgradeHandler).ServeHTTP(respRec, req)

	if respRec.Result().StatusCode != http.StatusNoContent {
		t.Errorf("Unexpected HTTP error code %d != %d", respRec.Result().StatusCode, http.StatusNoContent)
	}

	s, _, _ = cc_fb.LoadGameState(getGameIDFromPath(g.Path))
	if !s.GameData.UpgradeStatus[cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("Game state does not reflect upgrade bought")
	}
}

func init() {
	flag.Parse()

	cc_fb_config.SetCCFirebaseConfig()
	if cc_fb_config.CC_FIREBASE_CONFIG.Environment != cc_fb_config.DEV {
		panic(fmt.Sprintf("Firebase environment is not %s", cc_fb_config.ENVIRONMENT_TYPE_LOOKUP[cc_fb_config.DEV]))
	}
}
