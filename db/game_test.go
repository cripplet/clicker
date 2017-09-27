package cc_fb

import (
	"flag"
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/lib"
	"reflect"
	"testing"
)

func TestToFromInternalFBGameObservableData(t *testing.T) {
	cookiesPerClick := 1.0
	cps := 2.0
	buildingCost := 10.0
	upgradeCost := 12.0

	gameObservables := FBGameObservableData{
		CookiesPerClick: cookiesPerClick,
		CPS:             cps,
		BuildingCost: map[cookie_clicker.BuildingType]float64{
			cookie_clicker.BUILDING_TYPE_MOUSE: buildingCost,
		},
		UpgradeCost: map[cookie_clicker.UpgradeID]float64{
			cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER: upgradeCost,
		},
	}

	internalGameObservables := internalFBGameObservableData{
		CookiesPerClick: cookiesPerClick,
		CPS:             cps,
		BuildingCost: map[string]float64{
			cookie_clicker.BUILDING_TYPE_LOOKUP[cookie_clicker.BUILDING_TYPE_MOUSE]: buildingCost,
		},
		UpgradeCost: map[string]float64{
			cookie_clicker.UPGRADE_ID_LOOKUP[cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER]: upgradeCost,
		},
	}

	if !reflect.DeepEqual(toInternalFBGameObservableData(gameObservables), internalGameObservables) {
		t.Errorf("Error converting to internal game observables")
	}
	if !reflect.DeepEqual(fromInternalFBGameObservableData(internalGameObservables), gameObservables) {
		t.Errorf("Error converting from internal game observables")
	}
}

func TestToFromInternalFBGameState(t *testing.T) {
	version := "some-version"
	nCookies := 10.0
	nMice := 2
	upgradeBought := true
	gameID := "some-id"
	exist := true
	cookiesPerClick := 1.0
	cps := 2.0
	buildingCost := 10.0
	upgradeCost := 12.0

	gameObservables := FBGameObservableData{
		CookiesPerClick: cookiesPerClick,
		CPS:             cps,
		BuildingCost: map[cookie_clicker.BuildingType]float64{
			cookie_clicker.BUILDING_TYPE_MOUSE: buildingCost,
		},
		UpgradeCost: map[cookie_clicker.UpgradeID]float64{
			cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER: upgradeCost,
		},
	}
	gameData := cookie_clicker.GameStateData{
		Version:  version,
		NCookies: nCookies,
		NBuildings: map[cookie_clicker.BuildingType]int{
			cookie_clicker.BUILDING_TYPE_MOUSE: nMice,
		},
		UpgradeStatus: map[cookie_clicker.UpgradeID]bool{
			cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER: upgradeBought,
		},
	}
	gameState := FBGameState{
		ID:              gameID,
		Exist:           exist,
		GameData:        gameData,
		GameObservables: gameObservables,
	}

	internalGameObservables := internalFBGameObservableData{
		CookiesPerClick: cookiesPerClick,
		CPS:             cps,
		BuildingCost: map[string]float64{
			cookie_clicker.BUILDING_TYPE_LOOKUP[cookie_clicker.BUILDING_TYPE_MOUSE]: buildingCost,
		},
		UpgradeCost: map[string]float64{
			cookie_clicker.UPGRADE_ID_LOOKUP[cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER]: upgradeCost,
		},
	}
	internalGameData := internalGameStateData{
		Version:  version,
		NCookies: nCookies,
		NBuildings: map[string]int{
			cookie_clicker.BUILDING_TYPE_LOOKUP[cookie_clicker.BUILDING_TYPE_MOUSE]: nMice,
		},
		UpgradeStatus: map[string]bool{
			cookie_clicker.UPGRADE_ID_LOOKUP[cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER]: upgradeBought,
		},
	}
	internalGameState := internalFBGameState{
		ID:              gameID,
		Exist:           exist,
		GameData:        internalGameData,
		GameObservables: internalGameObservables,
	}

	if !reflect.DeepEqual(toInternalFBGameState(gameState), internalGameState) {
		t.Error("Error converting to internal game state")
	}
	if !reflect.DeepEqual(fromInternalFBGameState(internalGameState), gameState) {
		t.Error("Error converting from internal game state")
	}
}

func TestNewGame(t *testing.T) {
	ResetEnvironment(t)
	g, err := newGameState()
	if err != nil {
		t.Errorf("Unexpected error when loading game state: %v", err)
	}

	if g.ID == "" {
		t.Errorf("Game ID was not set")
	}
}

func TestLoadGame(t *testing.T) {
	ResetEnvironment(t)
	g, _ := newGameState()

	h, _, err := LoadGameState(g.ID)
	if err != nil {
		t.Errorf("Unexpected error when resetting database: %v", err)
	}

	if h.ID != g.ID {
		t.Errorf("Loaded game ID does not match: %s != %s", h.ID, g.ID)
	}
}

func TestLoadNonexistentGame(t *testing.T) {
	ResetEnvironment(t)
	g, _, err := LoadGameState("some-id")

	if err != nil {
		t.Errorf("Unexpected error when resetting database: %v", err)
	}

	if g.ID != "" {
		t.Errorf("Found game with given ID: %s", g.ID)
	}
}

func TestSaveGame(t *testing.T) {
	ResetEnvironment(t)
	g, _ := newGameState()
	g.ID = "some-other-id"

	err := SaveGameState(g, "null_etag")
	if err != nil {
		t.Errorf("Unexpected error when saving game state: %v", err)
	}

	h, _, _ := LoadGameState(g.ID)
	if !h.Exist {
		t.Errorf("Could not find game %s", g.ID)
	}
}

func init() {
	flag.Parse()

	cc_fb_config.SetCCFirebaseConfig()
	if cc_fb_config.CC_FIREBASE_CONFIG.Environment != cc_fb_config.DEV {
		panic(fmt.Sprintf("Firebase environment is not %s", cc_fb_config.ENVIRONMENT_TYPE_LOOKUP[cc_fb_config.DEV]))
	}
}
