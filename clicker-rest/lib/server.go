package cc_rest_lib

import (
	"encoding/json"
	"fmt"
	"github.com/cripplet/clicker/db"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/lib"
	"io/ioutil"
	"net/http"
	"regexp"
)

var newGameRegex *regexp.Regexp = regexp.MustCompile(`^/game(/)?$`)
var clickRegex *regexp.Regexp = regexp.MustCompile(`^/game/(?P<gameID>[\w-]*)/cookie/click(/)?$`)

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
	Hash   string `json:"hash"`
}

// TODO(cripplet): Implement multi-click hashing.
func ClickHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(clickRegex, req.URL.Path)["gameID"]
	switch {
	case req.Method == http.MethodPost:
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

		g := cookie_clicker.NewGameState()
		g.Load(s.GameData)
		for i := 0; i < clickRequest.NTimes; i++ {
			g.Click()
		}
		s.GameData = g.Dump()
		cc_fb.SaveGameState(s, eTag)
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
	}
}
