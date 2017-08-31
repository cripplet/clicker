package cookie_clicker

import (
	"fmt"
	"testing"
)

func TestMakeSimpleUpgrade(t *testing.T) {
	s := NewGameState()
	u := NewSimpleBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		2,
	)
	if (*u).GetName() != "New Upgrade" {
		t.Error(fmt.Sprintf("Expected name %s, got %s", "New Upgrade", (*u).GetName()))
	}
	if (*u).GetBuildingMultiplier(s) != 2 {
		t.Error(fmt.Sprintf("Expected BMul %e, got %e", 2, (*u).GetBuildingMultiplier(s)))
	}
}
