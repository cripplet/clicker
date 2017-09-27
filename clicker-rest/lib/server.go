// Package cc_rest_lib provides the handlers for managing and mutating CookieClicker games to the Firebase DB.
package cc_rest_lib

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/cripplet/clicker/db"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/lib"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var NewGameRegex *regexp.Regexp = regexp.MustCompile(`^/game(/)?$`)
var ClickRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/cookie/click(/)?$`)
var MineRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/cookie/mine(/)?$`)
var BuildingRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/building/(?P<buildingType>[\w-]*)(/)?$`)
var UpgradeRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/upgrade/(?P<upgradeID>[\w-]*)(/)?$`)

func regexpMatchNamedGroups(r *regexp.Regexp, s string) map[string]string {
	ret := map[string]string{}
	match := r.FindStringSubmatch(s)
	for index, name := range r.SubexpNames() {
		if name != "" {
			ret[name] = match[index]
		}
	}
	return ret
}

// NewGameResponse encapsulates the game state resource path.
type NewGameResponse struct {
	// The full URL of the game state,
	// e.g. https://some-project-id.firebaseio.com/game/some-game-id.json
	Path string `json:"path"`
}

// NewGameHandler creates a new game and stores the game state in the database.
// The URL must match NewGameRegex. The only valid method is a POST request
// with an empty body. A successful POST request will return a 201 CREATED
// response with a supplied NewGameResponse.
func NewGameHandler(resp http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodPost:
		s, _, err := cc_fb.LoadGameState("")
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		g := NewGameResponse{
			Path: fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, s.ID),
		}
		gJSON, err := json.Marshal(&g)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Header().Add("Content-Type", "application/json")
		resp.Write(gJSON)
		resp.WriteHeader(http.StatusCreated)
		break
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type ClickEvent struct {
	// Byte array of the correct validation SHA256 hash
	// associated with clicking the cookie NTimes.
	//
	// JSON will represent this as a Base64-encoded string,
	// e.g. LUt2MTZEOXc3d3I2dVV0RExFS2Q=.
	//
	// The current hash is provided in the game state. The correct Hash
	// value can be calculated by hashing (unencoded byte array) NTimes.
	Hash []byte `json:"hash"`

	// Time that this click event occured.
	ClickTime time.Time `json:"click_time"`
}

type ClickRequest struct {
	// List of click events, where each each event's ClickTime is
	// monotonically increasing
	// (i.e. Clicks[i].ClickTime < Clicks[i + 1].ClickTime)
	Clicks []ClickEvent `json:"clicks"`
}

// ClickHandler mutates the given game by incrementing the number of cookies
// in the bank by clicking "the cookie" some number of times. The URL must
// match ClickHandlerRegex. The only valid method is a POST request with a
// supplied ClickRequest body. A successful POST request will return a
// 204 NO CONTENT response with an empty body. ClickHandler may return with
// 412 PRECONDITION FAILED in the case of a race condition. The caller is
// responsible for reissuing the request.
func ClickHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(ClickRegex, req.URL.Path)["gameID"]
	switch req.Method {
	case http.MethodPost:
		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			break
		}

		fmt.Printf("%s", string(content))
		clickRequest := ClickRequest{}
		err = json.Unmarshal(content, &clickRequest)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			break
		}

		s, eTag, err := cc_fb.LoadGameState(gameID)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			break
		}

		if !s.Exist {
			resp.WriteHeader(http.StatusNotFound)
			break
		}

		g := cookie_clicker.NewGameState()
		g.Load(s.GameData)

		hash := s.Metadata.ClickHash
		clickTime := s.Metadata.ClickTime
		valid := true

		for _, clickEvent := range clickRequest.Clicks {
			newHashBytes := sha256.Sum256(hash)
			hash = newHashBytes[:]
			valid = valid && bytes.Equal(hash, clickEvent.Hash) && clickTime.Before(clickEvent.ClickTime)
			if valid {
				g.Click()
			} else {
				break
			}
			clickTime = clickEvent.ClickTime
		}
		if !valid {
			resp.WriteHeader(http.StatusBadRequest)
			break
		}

		s.Metadata.ClickHash = hash
		s.Metadata.ClickTime = clickTime

		s.GameData = g.Dump()

		err = cc_fb.SaveGameState(s, eTag)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusPreconditionFailed)
			break
		}

		resp.WriteHeader(http.StatusNoContent)
		break
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// MineHandler mutates the given game by incrementing the number of cookies
// in the bank by a calculated amount scaled since the last time this method
// was called. The URL must match MineHandlerRegex. The only valid method is
// a POST request with an empty body. A successful POST request will return a
// 204 NO CONTENT response with an empty body. MineHandler may return with
// 412 PRECONDITION FAILED in the case of a race condition. The caller is
// responsible for reissuing the request.
func MineHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(MineRegex, req.URL.Path)["gameID"]
	switch req.Method {
	case http.MethodPost:
		s, eTag, err := cc_fb.LoadGameState(gameID)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			break
		}

		if !s.Exist {
			resp.WriteHeader(http.StatusNotFound)
			break
		}

		g := cookie_clicker.NewGameState()
		g.Load(s.GameData)
		g.MineCookies(s.Metadata.MineTime, time.Now())
		s.GameData = g.Dump()

		err = cc_fb.SaveGameState(s, eTag)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusPreconditionFailed)
			break
		}

		resp.WriteHeader(http.StatusNoContent)
		break
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// BuildingHandler mutates the given game by incrementing the given building
// by spending cookies. The URL must match BuildingRegex. The only valid
// method is a POST request with an empty body. A successful POST request will
// return a 204 NO CONTENT response with an empty body. BuildingHandler will
// return with 402 PAYMENT REQUIRED in the case the user does not have enough
// cookies to buy a new building. BuildingHandler may return with
// 412 PRECONDITION FAILED in the case of a race condition. The caller is
// responsible for reissuing the request.
func BuildingHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(BuildingRegex, req.URL.Path)["gameID"]
	buildingType, present := cookie_clicker.BUILDING_TYPE_REVERSE_LOOKUP[regexpMatchNamedGroups(BuildingRegex, req.URL.Path)["buildingType"]]
	if !present {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	switch req.Method {
	case http.MethodPost:
		s, eTag, err := cc_fb.LoadGameState(gameID)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			break
		}

		if !s.Exist {
			resp.WriteHeader(http.StatusNotFound)
			break
		}

		g := cookie_clicker.NewGameState()
		g.Load(s.GameData)
		if !g.BuyBuilding(buildingType) {
			resp.WriteHeader(http.StatusPaymentRequired)
			break
		}
		s.GameData = g.Dump()
		s.GameObservables = cc_fb.GenerateFBGameObservableData(g)

		err = cc_fb.SaveGameState(s, eTag)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusPreconditionFailed)
			break
		}

		resp.WriteHeader(http.StatusNoContent)
		break
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// UpgradeHandler mutates the given game by incrementing the given upgrade by
// spending cookies. The URL must match UpgradeRegex. The only valid method is
// a POST request with an empty body. A successful POST request will return a
// 204 NO CONTENT response with an empty body. UpgradeHandler will
// return with 402 PAYMENT REQUIRED in the case the user does not have enough
// cookies to buy the upgrade. UpgradeHandler will return with
// 412 PRECONDITION FAILED in the case of a race condition. The caller is
// responsible for reissuing the request.
func UpgradeHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(UpgradeRegex, req.URL.Path)["gameID"]
	upgradeID, present := cookie_clicker.UPGRADE_ID_REVERSE_LOOKUP[regexpMatchNamedGroups(UpgradeRegex, req.URL.Path)["upgradeID"]]
	if !present {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	switch req.Method {
	case http.MethodPost:
		s, eTag, err := cc_fb.LoadGameState(gameID)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			break
		}

		if !s.Exist {
			resp.WriteHeader(http.StatusNotFound)
			break
		}

		g := cookie_clicker.NewGameState()
		g.Load(s.GameData)
		if !g.BuyUpgrade(upgradeID) {
			resp.WriteHeader(http.StatusPaymentRequired)
			break
		}
		s.GameData = g.Dump()
		s.GameObservables = cc_fb.GenerateFBGameObservableData(g)

		err = cc_fb.SaveGameState(s, eTag)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusPreconditionFailed)
			break
		}

		resp.WriteHeader(http.StatusNoContent)
		break
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GameRouter is responsible for routing all URLs to an appropriate handler.
func GameRouter(resp http.ResponseWriter, req *http.Request) {
	switch {
	case NewGameRegex.MatchString(req.URL.Path):
		NewGameHandler(resp, req)
		break
	case ClickRegex.MatchString(req.URL.Path):
		ClickHandler(resp, req)
		break
	case MineRegex.MatchString(req.URL.Path):
		MineHandler(resp, req)
		break
	case BuildingRegex.MatchString(req.URL.Path):
		BuildingHandler(resp, req)
		break
	case UpgradeRegex.MatchString(req.URL.Path):
		UpgradeHandler(resp, req)
		break
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}
