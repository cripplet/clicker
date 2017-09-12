package cc_fb

import (
	"fmt"
	"github.com/cripplet/clicker/db/firebase/config"
	"github.com/cripplet/clicker/firebase-db"
	"net/http"
	"testing"
)

func ResetEnvironment(t *testing.T) {
	_, status_code, err := firebase_db.Delete(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath),
		false,
		"",
		map[string]string{},
	)
	if err != nil {
		t.Errorf("Unexpected error when resetting database: %v", err)
	}
	if status_code != http.StatusOK {
		t.Errorf("Unexpected HTTP status code when deleting root directory: %d != %d", status_code, http.StatusOK)
	}
}

func TestLoadGame(t *testing.T) {
	ResetEnvironment(t)
	g, u, _ := CreateGameState()

	if g.ID != u.GameID {
		t.Errorf("Game ID does not match user token: %s != %s", g.ID, u.GameID)
	}
}

func TestLoadGameStateBadAuthorization(t *testing.T) {
	ResetEnvironment(t)
	id := "some-id"
	_, _, err := LoadGameState(id, "some-incorrect-token")
	if err == nil {
		t.Error("No error was returned when expecting one.")
	}
}

func init() {
	if cc_fb_config.CC_FIREBASE_CONFIG.Environment != cc_fb_config.DEV {
		panic(fmt.Sprintf("Firebase environment is not %s", cc_fb_config.ENVIRONMENT_TYPE_LOOKUP[cc_fb_config.DEV]))
	}
}
