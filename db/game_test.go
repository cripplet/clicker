package cc_fb

import (
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/lib"
	"reflect"
	"testing"
)

func TestToFromInternalFBGameState(t *testing.T) {
	version := "some-version"
	nCookies := 10.0
	nMice := 2
	upgradeBought := true
	gameID := "some-id"
	exist := true

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
		ID:       gameID,
		Exist:    exist,
		GameData: gameData,
	}

	internalGameData := internalGameStateData{
		Version:  version,
		NCookies: nCookies,
		NBuildings: map[string]int{
			fmt.Sprintf("_%d", cookie_clicker.BUILDING_TYPE_MOUSE): nMice,
		},
		UpgradeStatus: map[string]bool{
			fmt.Sprintf("_%d", cookie_clicker.UPGRADE_ID_REINFORCED_INDEX_FINGER): upgradeBought,
		},
	}
	internalGameState := internalFBGameState{
		ID:       gameID,
		Exist:    exist,
		GameData: internalGameData,
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
	g, err := NewGameState()
	if err != nil {
		t.Errorf("Unexpected error when loading game state: %v", err)
	}

	if g.ID == "" {
		t.Errorf("Game ID was not set")
	}
}

func TestLoadGame(t *testing.T) {
	ResetEnvironment(t)
	g, _ := NewGameState()

	h, err := LoadGameState(g.ID)
	if err != nil {
		t.Errorf("Unexpected error when resetting database: %v", err)
	}

	if h.ID != g.ID {
		t.Errorf("Loaded game ID does not match: %s != %s", h.ID, g.ID)
	}
}

func TestLoadNonexistentGame(t *testing.T) {
	ResetEnvironment(t)
	g, err := LoadGameState("some-id")

	if err != nil {
		t.Errorf("Unexpected error when resetting database: %v", err)
	}

	if g.ID != "" {
		t.Errorf("Found game with given ID: %s", g.ID)
	}
}

func TestSaveGame(t *testing.T) {
	ResetEnvironment(t)
	g, _ := NewGameState()
	g.ID = "some-other-id"

	err := SaveGameState(g)
	if err != nil {
		t.Errorf("Unexpected error when saving game state: %v", err)
	}

	h, _ := LoadGameState(g.ID)
	if !h.Exist {
		t.Errorf("Could not find game %s", g.ID)
	}
}

func init() {
	if cc_fb_config.CC_FIREBASE_CONFIG.Environment != cc_fb_config.DEV {
		panic(fmt.Sprintf("Firebase environment is not %s", cc_fb_config.ENVIRONMENT_TYPE_LOOKUP[cc_fb_config.DEV]))
	}
}
