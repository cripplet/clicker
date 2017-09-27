// Package cc_fb handles CookieClicker game state representation in the Firebase DB.
// TODO(cripplet): Move to github.com/cripplet/clicker-rest/db.
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

// FBGameState is the exported state representation of Firebase database
// data. Note that this is not the "true" JSON output as seen when pulling
// from the database. The JSON field names are included for convenience, but
// this type is never actually serialized.
type FBGameState struct {
	// Unique game ID.
	ID string `json:"id"`

	// Placeholder boolean; when loading a nonexistent game JSON, this
	// field will be the boolean zero value, i.e. false. This field is
	// explictly set to true for all games.
	Exist bool `json:"exist"`

	// Imported game data, used to initialize a game. Due to Firebase
	// shenanigans, this struct's map[ENUM]VALUE properties are instead
	// set to map[STRING]VALUE instead. See FBGameObservableData for
	// related comments.
	GameData cookie_clicker.GameStateData `json:"data"`

	// Calculated data from a game. These fields are read-only fields
	// and will be automatically updated when the game has mutated.
	GameObservables FBGameObservableData `json:"observables"`

	// Metadata kept by the REST server.
	Metadata FBGameMetadata `json:"metadata"`
}

// FBGameMetadata is the metadata specific to the REST server. The actual game
// has no knowledge of this data.
type FBGameMetadata struct {
	// Data useful for discouraging spamming large amount of REST click
	// requests.
	ClickHash []byte `json:"click_hash"`

	// Last time that cookies added to the game due to CPS contributions
	// have been calculated.
	MineTime time.Time `json:"mine_time"`
}

// FBGameObservableData is calculated from the game and is read-only.
type FBGameObservableData struct {
	// Number of cookies added to the game per click.
	CookiesPerClick float64 `json:"cookies_per_click"`

	// Number of cookies added to the game per second.
	CPS float64 `json:"cps"`

	// The number of cookies necessary to buy a building of that type.
	// The key is translated from the cookie_clicker.BuildingType enum, via
	// cookie_clicker.BUILDING_TYPE_LOOKUP.
	BuildingCost map[string]float64 `json:"building_cost"`

	// The number of cookies necessary to buy the given upgrade. The key is
	// translated via cookie_clicker.UPGRADE_ID_LOOKUP.
	UpgradeCost map[string]float64 `json:"upgrade_cost"`
}

type internalGameStateData struct {
	Version       string          `json:"version"`
	NCookies      float64         `json:"n_cookies"`
	NBuildings    map[string]int  `json:"buildings"`
	UpgradeStatus map[string]bool `json:"upgrades"`
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

type postID struct {
	Name string `json:"name"`
}

// GenerateFBGameObservableData will construct a representation of the
// read-only game fields.
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

	p := postID{}
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

// LoadGameState retrieves a game from the Firebase DB (including all related
// metadata). If the supplied ID is empty, construct a new game instead (and
// commit to the DB). Returns an ETag value alongside the game, which must be
// passed into SaveGameState. A load / save cycle is atomic (or will fail).
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

// SaveGameState commits a game to the Firebase DB. ETag was supplied
// previously when loading the game. LoadGameState will fail if the given
// ETag does not match the current DB ETag. The user is responsible for
// retrying the load / save cycle in case of failure.
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
