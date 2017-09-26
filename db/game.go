package cc_fb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/firebase-db"
	"github.com/cripplet/clicker/lib"
	"net/http"
	"time"
)

type FBGameState struct {
	ID    string
	Exist bool

	// Imported Game data, used to initialize a game.
	GameData cookie_clicker.GameStateData

	// Calculated Game data.
	GameObservables FBGameObservableData

	// Metadata kept by the REST server.
	Metadata FBGameMetadata
}

type FBGameMetadata struct {
	ClickHash []byte    `json:"click_hash"`
	MineTime  time.Time `json:"mine_time"`
}

type FBGameObservableData struct {
	CookiesPerClick float64            `json:"cookies_per_click"`
	CPS             float64            `json:"cps"`
	BuildingCost    map[string]float64 `json:"building_cost"`
	UpgradeCost     map[string]float64 `json:"upgrade_cost"`
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
	Metadata        FBGameMetadata        `json:"metadata"`
}

func toInternalFBGameState(s FBGameState) internalFBGameState {
	internalData := internalGameStateData{
		Version:       s.GameData.Version,
		NCookies:      s.GameData.NCookies,
		NBuildings:    make(map[string]int),
		UpgradeStatus: make(map[string]bool),
	}
	for buildingType, nBuildings := range s.GameData.NBuildings {
		internalData.NBuildings[cookie_clicker.BUILDING_TYPE_LOOKUP[buildingType]] = nBuildings
	}
	for upgradeID, bought := range s.GameData.UpgradeStatus {
		internalData.UpgradeStatus[cookie_clicker.UPGRADE_ID_LOOKUP[upgradeID]] = bought
	}
	return internalFBGameState{
		ID:              s.ID,
		Exist:           s.Exist,
		GameData:        internalData,
		GameObservables: s.GameObservables,
		Metadata:        s.Metadata,
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
		gameData.NBuildings[cookie_clicker.BUILDING_TYPE_REVERSE_LOOKUP[buildingTypeString]] = nBuildings
	}
	for upgradeIDString, bought := range s.GameData.UpgradeStatus {
		gameData.UpgradeStatus[cookie_clicker.UPGRADE_ID_REVERSE_LOOKUP[upgradeIDString]] = bought
	}
	return FBGameState{
		ID:              s.ID,
		Exist:           s.Exist,
		GameData:        gameData,
		GameObservables: s.GameObservables,
		Metadata:        s.Metadata,
	}
}

type PostID struct {
	Name string `json:"name"`
}

func GenerateFBGameObservableData(g *cookie_clicker.GameStateStruct) FBGameObservableData {
	n := time.Now()

	buildingCost := map[string]float64{}
	for buildingType, building := range g.GetBuildings() {
		buildingCost[cookie_clicker.BUILDING_TYPE_LOOKUP[buildingType]] = building.GetCost(g.GetNBuildings()[buildingType] + 1)
	}
	upgradeCost := map[string]float64{}
	for upgradeID, upgrade := range g.GetUpgrades() {
		upgradeCost[cookie_clicker.UPGRADE_ID_LOOKUP[upgradeID]] = upgrade.GetCost(g)
	}
	o := FBGameObservableData{
		CookiesPerClick: g.GetCookiesPerClick(),
		CPS:             g.GetCPS(n, n.Add(time.Second)),
		BuildingCost:    buildingCost,
		UpgradeCost:     upgradeCost,
	}
	return o
}

func newGameState() (FBGameState, error) {
	d := cookie_clicker.NewGameStateData()
	g := cookie_clicker.NewGameState()
	g.Load(*d)

	s := FBGameState{
		Exist:           true,
		GameData:        *d,
		GameObservables: GenerateFBGameObservableData(g),
	}

	i := toInternalFBGameState(s)
	iJSON, err := json.Marshal(&i)
	if err != nil {
		return FBGameState{}, err
	}

	p := PostID{}
	_, _, eTag, err := firebase_db.Post(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s/game.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath),
		iJSON,
		true,
		map[string]string{},
		&p,
	)
	if err != nil {
		return FBGameState{}, err
	}

	i.ID = p.Name
	i.Metadata.ClickHash = []byte(i.ID)
	i.Metadata.MineTime = time.Now()

	err = SaveGameState(fromInternalFBGameState(i), eTag)
	if err != nil {
		return FBGameState{}, err
	}
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
