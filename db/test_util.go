package cc_fb

import (
	"fmt"
	"github.com/cripplet/clicker/db/config"
	"github.com/cripplet/clicker/firebase-db"
	"net/http"
	"testing"
)

func ResetEnvironment(t *testing.T) {
	_, statusCode, _, err := firebase_db.Delete(
		cc_fb_config.CC_FIREBASE_CONFIG.Client,
		fmt.Sprintf("%s.json", cc_fb_config.CC_FIREBASE_CONFIG.ProjectPath),
		false,
		"",
		map[string]string{},
	)
	if err != nil {
		t.Errorf("Unexpected error when resetting database: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP status code when deleting root directory: %d != %d", statusCode, http.StatusOK)
	}
}
