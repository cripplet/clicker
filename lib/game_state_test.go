package cookie_clicker

import (
	"fmt"
	"testing"
	"time"
)

var start time.Time = time.Now()
var end time.Time = start.Add(time.Second)

func TestMakeGameState(t *testing.T) {
	s := NewGameState()

	var i BuildingType
	for i = 0; i < BUILDING_TYPE_ENUM_EOF; i++ {
		if s.nBuildings[i] != 0 {
			t.Error(fmt.Sprintf("Expected %d instances of type %d, got %d", 0, i, s.nBuildings[i]))
		}
	}
}

func TestAddCookies(t *testing.T) {
	s := NewGameState()
	s.addCookies(1)
	if s.GetCookies() != 1 {
		t.Error(fmt.Sprintf("Expected %e cookies, got %e", 1, s.GetCookies()))
	}
}

func TestSubtractCookiesTooExpensive(t *testing.T) {
	s := NewGameState()
	if s.subtractCookies(1) {
		t.Error("Subtracted cookies when not possible")
	}
}

func TestSubtractCookies(t *testing.T) {
	s := NewGameState()
	s.addCookies(1)
	if !s.subtractCookies(1) {
		t.Error("Could not subtract cookies")
	}
	if s.GetCookies() != 0 {
		t.Error("Expected %e cookies, got %e", 0, s.GetCookies())
	}
}

func TestCalculateCPSNoUpgrades(t *testing.T) {
	s := NewGameState()

	s.loadBuildingCPSRef(BUILDING_CPS_LOOKUP)
	s.nBuildings[BUILDING_TYPE_MOUSE] = 1

	cps := s.calculateCPS()
	if cps != 0.2 {
		t.Error(fmt.Sprintf("Expected total CPS %e, got %e", 0.2, cps))
	}
}

func TestCalculateCPSSimpleUpgrade(t *testing.T) {
	s := NewGameState()

	s.loadBuildingCPSRef(BUILDING_CPS_LOOKUP)
	s.nBuildings[BUILDING_TYPE_MOUSE] = 1
	s.upgradeStatus[UPGRADE_ID_REINFORCED_INDEX_FINGER] = true

	u := NewBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		0,
		2,
		0,
	)

	s.loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	cps := s.calculateCPS()
	if s.calculateCPS() != 0.4 {
		t.Error(fmt.Sprintf("Expected total CPS %e, got %e", 0.4, cps))
	}
}

func TestBuyNonexistentUpgrade(t *testing.T) {
	s := NewGameState()

	if s.BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Bought nonexistent upgrade.")
	}
	if s.upgradeStatus[UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("Upgrade was bought even though BuyUpgrade returned not bought.")
	}
}

func TestDoubleBuyUpgrade(t *testing.T) {
	s := NewGameState()

	u := NewBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		100,
		2,
		0,
	)

	s.loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	s.upgradeStatus[UPGRADE_ID_REINFORCED_INDEX_FINGER] = true

	if s.BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Double bought upgrade.")
	}
	if !s.upgradeStatus[UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("Unbought upgrade.")
	}
}

func TestBuyUpgradeTooExpensive(t *testing.T) {
	s := NewGameState()

	u := NewBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		100,
		2,
		0,
	)

	s.loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if s.BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Bought too expensive an upgrade.")
	}
}

func TestBuyUpgradeLocked(t *testing.T) {
	s := NewGameState()

	u := NewBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		0,
		2,
		1,
	)

	s.loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if s.BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Bought a locked upgrade")
	}
}

func TestBuyCPSUpgrade(t *testing.T) {
	s := NewGameState()

	s.loadBuildingCPSRef(BUILDING_CPS_LOOKUP)
	s.nBuildings[BUILDING_TYPE_MOUSE] = 1

	u := NewBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		0,
		2,
		0,
	)

	s.loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if !s.BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Could not buy simple upgrade.")
	}

	if !s.upgradeStatus[UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("BuyUpgrade reported upgrade was bought, but does not show up as bought.")
	}

	if s.GetCPS(start, end) != 0.4 {
		t.Error(fmt.Sprintf("Expected CPS %e, got %e", 0.4, s.GetCPS(start, end)))
	}
}

func TestBuyCookiesPerClickUpgrade(t *testing.T) {
	s := NewGameState()

	s.loadCookiesPerClickRef(2)

	u := NewBasicClickUpgrade(
		"New Upgrade",
		0,
		3,
	)

	s.loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if !s.BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Could not buy click upgrade.")
	}

	if !s.upgradeStatus[UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("BuyUpgrade reported upgrade was bought, but does not show up as bought.")
	}

	if s.GetCookiesPerClick() != 6 {
		t.Error(fmt.Sprintf("Expected %e cookies per click, got %e", 6, s.GetCookiesPerClick()))
	}
}

func TestBuyBuildingFree(t *testing.T) {
	s := NewGameState()
	s.loadBuildingCPSRef(BUILDING_CPS_LOOKUP)
	s.BuyBuilding(BUILDING_TYPE_MOUSE)
	if s.GetCPS(start, end) != 0.2 {
		t.Error(fmt.Sprintf("Expected CPS %e, got %e", 0.2, s.GetCPS(start, end)))
	}
}

func TestBuyBuildingTooExpensive(t *testing.T) {
	s := NewGameState()
	s.loadBuildingCost(map[BuildingType]BuildingCostFunction{
		BUILDING_TYPE_MOUSE: func(current int) float64 { return 1 },
	})
	if s.BuyBuilding(BUILDING_TYPE_MOUSE) {
		t.Error("Expected to not buy building, but bought anyways.")
	}
	if s.nBuildings[BUILDING_TYPE_MOUSE] != 0 {
		t.Error(fmt.Sprintf("Expected %d buildings, got %d", 0, s.nBuildings[BUILDING_TYPE_MOUSE]))
	}
}

func TestBuyBuildingAffordable(t *testing.T) {
	s := NewGameState()
	s.addCookies(1)
	s.loadBuildingCost(map[BuildingType]BuildingCostFunction{
		BUILDING_TYPE_MOUSE: func(current int) float64 { return 1 },
	})
	if !s.BuyBuilding(BUILDING_TYPE_MOUSE) {
		t.Error("Expected to buy building, but couldn't.")
	}
	if s.nBuildings[BUILDING_TYPE_MOUSE] != 1 {
		t.Error(fmt.Sprintf("Expected %d buildings, got %d", 1, s.nBuildings[BUILDING_TYPE_MOUSE]))
	}
	if s.GetCookies() != 0 {
		t.Error(fmt.Sprintf("Expected %e cookies left, got %e", 0, s.GetCookies))
	}
}

func TestGameStateLoad(t *testing.T) {
	s := NewGameState()
	s.Load(GameStateData{})

	if s.buildingCPSRef[BUILDING_TYPE_MOUSE] != BUILDING_CPS_LOOKUP[BUILDING_TYPE_MOUSE] {
		t.Error("Could not load game state constants properly")
	}
}
