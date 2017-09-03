package cookie_clicker

import (
	"fmt"
	"testing"
)

func TestMakeBasicClickUpgrade(t *testing.T) {
	u := NewBasicClickUpgrade(
		"Simple Click",
		1,
		2,
	)

	if (*u).GetBuildingType() != BUILDING_TYPE_ENUM_EOF {
		t.Error(fmt.Sprintf("Expected building type %d, got %d", BUILDING_TYPE_ENUM_EOF, (*u).GetBuildingType))
	}

	s := NewGameState()
	if (*u).GetClickMultiplier(s) != 2 {
		t.Error(fmt.Sprintf("Expected click multiplier %e, got %e", 2, (*u).GetClickMultiplier(s)))
	}
}
