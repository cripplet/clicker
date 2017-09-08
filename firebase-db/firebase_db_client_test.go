package firebase_db

import (
	"testing"
)

func TestNewGoogleClient(t *testing.T) {
	_, err := NewGoogleClient(credentials)
	if err != nil {
		t.Errorf("Could not construct Firebase client, failed with error: %v", err)
	}
}
