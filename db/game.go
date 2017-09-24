package cc_fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/firebase-db"
	"github.com/cripplet/clicker/lib"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type FBGameState struct {
	ID       string
	Exist    bool
	GameData cookie_clicker.GameStateData
}

type internalGameStateData struct {
	Version       string          `json:"version"`
	NCookies      float64         `json:"n_cookies"`
	NBuildings    map[string]int  `json:"n_buildings"`
	UpgradeStatus map[string]bool `json:"upgrade_status"`
}

type internalFBGameState struct {
	ID       string                `json:"ID"`
	Exist    bool                  `json:"exist"`
	GameData internalGameStateData `json:"data"`
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
		ID:       s.ID,
		Exist:    s.Exist,
		GameData: internalData,
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
		ID:       s.ID,
		Exist:    s.Exist,
		GameData: gameData,
	}
}

func randomString(n int) string {
	r := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = r[rand.Intn(len(r))]
	}
	return string(b)
}

type PostID struct {
	Name string `json:"name"`
}

func NewGameState() (FBGameState, error) {
	d := cookie_clicker.NewGameStateData()
	g := toInternalFBGameState(FBGameState{
		Exist:    true,
		GameData: *d,
	})
	gJSON, err := json.Marshal(&g)
	if err != nil {
		return FBGameState{}, err
	}

	p := PostID{}

	_, _, err = firebase_db.Post(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath),
		gJSON,
		false,
		map[string]string{},
		&p,
	)
	if err != nil {
		return FBGameState{}, err
	}

	g.ID = p.Name
	return fromInternalFBGameState(g), nil
}

func LoadGameState(id string) (FBGameState, error) {
	if id == "" {
		return NewGameState()
	}

	g := internalFBGameState{}
	_, _, err := firebase_db.Get(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, id),
		false,
		map[string]string{},
		&g,
	)
	if err != nil {
		return FBGameState{}, err
	}
	if g.Exist {
		g.ID = id
	}
	return fromInternalFBGameState(g), nil
}

func SaveGameState(g FBGameState) error {
	i := toInternalFBGameState(g)
	iJSON, err := json.Marshal(&i)
	if err != nil {
		return err
	}

	_, statusCode, err := firebase_db.Put(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game/%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath, g.ID),
		iJSON,
		false,
		"",
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
