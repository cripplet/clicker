package cc_firebase

import (
	"flag"
	"fmt"
	cc_firebase_config "github.com/cripplet/clicker/db/firebase/config"
	"net/http"
	"testing"
)

var credentials string

func ResetDevEnvironment(t *testing.T) {
	c, _ := NewFirebaseClient(credentials)

	request, _ := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf(
			"%s/%s.json",
			cc_firebase_config.DB_CONFIG.BaseURL,
			ENVIRONMENT_LOOKUP[ENVIRONMENT_DEV],
		),
		nil,
	)
	_, err := c.Do(request)
	if err != nil {
		t.Errorf("Could not delete test DB, failed with error: %v", err)
	}
}

func TestNewFirebaseClient(t *testing.T) {
	_, err := NewFirebaseClient(credentials)
	if err != nil {
		t.Errorf("Could not construct Firebase client, failed with error: %v", err)
	}
}

func TestCreateSession(t *testing.T) {
	ResetDevEnvironment(t)
	c, _ := NewFirebaseClient(credentials)
	s, err := CreateSession(c, ENVIRONMENT_DEV, "some-id")
	if err != nil {
		t.Errorf("Could not create Session, failed with error: %v", err)
	}
	if s.ID != "some-id" {
		t.Errorf("Session ID mismatch, got %v", s.ID)
	}
}

func TestReadSessionDNE(t *testing.T) {
	ResetDevEnvironment(t)
	c, _ := NewFirebaseClient(credentials)
	s, _ := ReadSession(c, ENVIRONMENT_DEV, "non-existent")
	if s.ID != "" {
		t.Errorf("Found non-existent session %v", s.ID)
	}
}

func TestReadSession(t *testing.T) {
	ResetDevEnvironment(t)
	c, _ := NewFirebaseClient(credentials)
	s, _ := CreateSession(c, ENVIRONMENT_DEV, "some-id")
	sp, _ := ReadSession(c, ENVIRONMENT_DEV, s.ID)
	if sp.ID != s.ID {
		t.Errorf("Session ID mismatch: %v != %v", sp.ID, s.ID)
	}
}

func TestCloneSession(t *testing.T) {
	ResetDevEnvironment(t)
	c, _ := NewFirebaseClient(credentials)
	s, _ := CreateSession(c, ENVIRONMENT_DEV, "some-id")
	sp, _ := CloneSession(c, ENVIRONMENT_DEV, s.ID)
	if sp.ID != s.ID+"-clone" {
		t.Errorf("Cloned session ID did not match: %v != %v", s.ID+"-clone", sp.ID)
	}
}

func init() {
	flag.StringVar(&credentials, "credentials", "", "")
	flag.Parse()
}
