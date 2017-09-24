package cc_fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/firebase-db"
	"github.com/cripplet/clicker/lib"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FBGameState struct {
	ID    string
	Exist bool

	// Imported Game data, used to initialize a game.
	GameData cookie_clicker.GameStateData

	// Calculated Game data.
	GameObservables FBGameObservableData
}

type FBGameObservableData struct {
	CookiesPerClick float64 `json:"cookies_per_click"`
	CPS             float64 `json:"cps"`
}

type internalGameStateData struct {
	Version       string          `json:"version"`
	NCookies      float64         `json:"n_cookies"`
	NBuildings    map[string]int  `json:"building"`
	UpgradeStatus map[string]bool `json:"upgrade"`
}

type internalFBGameState struct {
	ID              string                `json:"id"`
	Exist           bool                  `json:"exist"`
	GameData        internalGameStateData `json:"data"`
	GameObservables FBGameObservableData  `json:"observables"`
}

func toInternalFBGameState(s FBGameState) internalFBGameState {
	internalData := internalGameStateData{
		Version:       s.GameData.Version,
		NCookies:      s.GameData.NCookies,
		NBuildings:    make(map[string]int),
		UpgradeStatus: make(map[string]bool),
	}
	for buildingType, nBuildings := range s.GameData.NBuildings {
		internalData.NBuildings[fmt.Sprintf("_%d", buildingType)] = nBuildings
	}
	for upgradeID, bought := range s.GameData.UpgradeStatus {
		internalData.UpgradeStatus[fmt.Sprintf("_%d", upgradeID)] = bought
	}
	return internalFBGameState{
		ID:              s.ID,
		Exist:           s.Exist,
		GameData:        internalData,
		GameObservables: s.GameObservables,
	}
}

func fromInternalFBGameState(s internalFBGameState) FBGameState {
	gameData := cookie_clicker.GameStateData{
		Version:       s.GameData.Version,
		NCookies:      s.GameData.NCookies,
		NBuildings:    make(map[cookie_clicker.BuildingType]int),
		UpgradeStatus: make(map[cookie_clicker.UpgradeID]bool),
	}
	for buildingTypeString, nBuildings := range s.GameData.NBuildings {
		buildingTypeInt, _ := strconv.Atoi(strings.Replace(buildingTypeString, "_", "", -1))
		gameData.NBuildings[cookie_clicker.BuildingType(buildingTypeInt)] = nBuildings
	}
	for upgradeIDString, bought := range s.GameData.UpgradeStatus {
		upgradeIDInt, _ := strconv.Atoi(strings.Replace(upgradeIDString, "_", "", -1))
		gameData.UpgradeStatus[cookie_clicker.UpgradeID(upgradeIDInt)] = bought
	}
	return FBGameState{
		ID:              s.ID,
		Exist:           s.Exist,
		GameData:        gameData,
		GameObservables: s.GameObservables,
	}
}

type PostID struct {
	Name string `json:"name"`
}

func newGameState() (FBGameState, error) {
	n := time.Now()

	d := cookie_clicker.NewGameStateData()
	g := cookie_clicker.NewGameState()
	g.Load(*d)
	o := FBGameObservableData{
		CookiesPerClick: g.GetCookiesPerClick(),
		CPS:             g.GetCPS(n, n.Add(time.Second)),
	}

	s := FBGameState{
		Exist:           true,
		GameData:        *d,
		GameObservables: o,
	}

	i := toInternalFBGameState(s)
	iJSON, err := json.Marshal(&i)
	if err != nil {
		return FBGameState{}, err
	}

	p := PostID{}
	_, _, _, err = firebase_db.Post(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath),
		iJSON,
		false,
		map[string]string{},
		&p,
	)
	if err != nil {
		return FBGameState{}, err
	}

	i.ID = p.Name
	return fromInternalFBGameState(i), nil
}

func LoadGameState(id string) (FBGameState, string, error) {
	if id == "" {
		g, err := newGameState()
		if err != nil {
			return FBGameState{}, "", nil
		}
		id = g.ID
	}

	i := internalFBGameState{}
	_, _, eTag, err := firebase_db.Get(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, id),
		true,
		map[string]string{},
		&i,
	)
	if err != nil {
		return FBGameState{}, "", err
	}
	if i.Exist {
		i.ID = id
	}
	return fromInternalFBGameState(i), eTag, nil
}

func SaveGameState(g FBGameState, eTag string) error {
	i := toInternalFBGameState(g)
	iJSON, err := json.Marshal(&i)
	if err != nil {
		return err
	}

	_, statusCode, _, err := firebase_db.Put(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, g.ID),
		iJSON,
		false,
		eTag,
		map[string]string{},
		nil,
	)

	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("HTTP error %d", statusCode))
	}

	return nil
}
