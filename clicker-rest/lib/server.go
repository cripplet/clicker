package cc_rest_lib

import (
	"encoding/json"
	"github.com/cripplet/clicker/db"
	"github.com/cripplet/clicker/lib"
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

func NewGameHandler(resp http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodPost:
		s, err := cc_fb.NewGameState()
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		sJSON, err := json.Marshal(s)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Write(sJSON)
		break
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// TODO(cripplet): implement the hash
// TODO(cripplet): implement etag
func ClickHandler(resp http.ResponseWriter, req *http.Request) {
	gameID := regexpMatchNamedGroups(clickRegex, req.URL.Path)["gameID"]
	switch {
	case req.Method == http.MethodPost:
		s, err := cc_fb.LoadGameState(gameID)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		if !s.Exist {
			resp.WriteHeader(http.StatusNotFound)
			return
		}
		g := cookie_clicker.NewGameState()
		g.Load(s.GameData)
		g.Click()
		s.GameData = g.Dump()
		cc_fb.SaveGameState(s)
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
