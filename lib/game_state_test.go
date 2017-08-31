package cookie_clicker

import (
	"fmt"
	"testing"
)

func TestMakeGameState(t *testing.T) {
	s := NewGameState()

	var i BuildingType
	for i = 0; i < BUILDING_TYPE_ENUM_EOF; i++ {
		if (*s).n_buildings[i] != 0 {
			t.Error(fmt.Sprintf("Expected %d instances of type %d, got %d", 0, i, (*s).n_buildings[i]))
		}
	}
}

func TestCalculateCPSNoUpgrades(t *testing.T) {
	s := NewGameState()
	s.n_buildings[BUILDING_TYPE_MOUSE] = 1

	upgrades := make(map[UpgradeID]UpgradeInterface)
	if (*s).CalculateCPS(upgrades) != 0.2 {
		t.Error(fmt.Sprintf("Expected total CPS %e, got %e", 0.2, (*s).cps))
	}
}

func TestCalculateCPSSimpleUpgrade(t *testing.T) {
	s := NewGameState()
	s.n_buildings[BUILDING_TYPE_MOUSE] = 1
	s.upgrade_status[UPGRADE_ID_REINFORCED_INDEX_FINGER] = true

	u := NewSimpleBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		2,
	)

	upgrades := map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	}

	if (*s).CalculateCPS(upgrades) != 0.4 {
		t.Error(fmt.Sprintf("Expected total CPS %e, got %e", 0.4, (*s).cps))
	}
}
