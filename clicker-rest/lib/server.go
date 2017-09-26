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

var newGameRegex *regexp.Regexp = regexp.MustCompile(`^/game(/)?$`)
var clickRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/cookie/click(/)?$`)
var mineRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/cookie/mine(/)?$`)
var buyBuildingRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/building/(?P<buildingType>[\w-]*)(/)?$`)
var buyUpgradeRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/upgrade/(?P<upgradeID>[\w-]*)(/)?$`)

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

type NewGameResponse struct {
	Path string `json:"path"`
}

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

type ClickRequest struct {
	NTimes int    `json:"n_times"`
	Hash   []byte `json:"hash"`
}

func ClickHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(clickRegex, req.URL.Path)["gameID"]
	switch req.Method {
	case http.MethodPost:
		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			break
		}

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

		hash := s.Metadata.ClickHash
		for i := 0; i < clickRequest.NTimes; i++ {
			newHashBytes := sha256.Sum256(hash)
			hash = newHashBytes[:]
		}

		if !bytes.Equal(clickRequest.Hash, hash) {
			http.Error(resp, "Invalid hash value provided", http.StatusBadRequest)
			break
		}

		s.Metadata.ClickHash = hash

		g := cookie_clicker.NewGameState()
		g.Load(s.GameData)
		for i := 0; i < clickRequest.NTimes; i++ {
			g.Click()
		}
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

func MineHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(mineRegex, req.URL.Path)["gameID"]
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

func BuyBuildingHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(buyBuildingRegex, req.URL.Path)["gameID"]
	buildingType, present := cookie_clicker.BUILDING_TYPE_REVERSE_LOOKUP[regexpMatchNamedGroups(buyBuildingRegex, req.URL.Path)["buildingType"]]
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

func BuyUpgradeHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(buyUpgradeRegex, req.URL.Path)["gameID"]
	upgradeID, present := cookie_clicker.UPGRADE_ID_REVERSE_LOOKUP[regexpMatchNamedGroups(buyUpgradeRegex, req.URL.Path)["upgradeID"]]
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

func GameRouter(resp http.ResponseWriter, req *http.Request) {
	switch {
	case newGameRegex.MatchString(req.URL.Path):
		NewGameHandler(resp, req)
		break
	case clickRegex.MatchString(req.URL.Path):
		ClickHandler(resp, req)
		break
	case mineRegex.MatchString(req.URL.Path):
		MineHandler(resp, req)
		break
	case buyBuildingRegex.MatchString(req.URL.Path):
		BuyBuildingHandler(resp, req)
		break
	case buyUpgradeRegex.MatchString(req.URL.Path):
		BuyUpgradeHandler(resp, req)
		break
	}
}
