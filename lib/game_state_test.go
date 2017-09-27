package cookie_clicker

import (
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
			t.Errorf("Expected %d instances of type %d, got %d", 0, i, s.nBuildings[i])
		}
	}

	var u UpgradeID
	for u = 0; u < UPGRADE_ID_ENUM_EOF; u++ {
		if s.upgradeStatus[u] {
			t.Errorf("Upgrade %d was bought when it shouldn't have been", u)
		}
	}
}

func TestAddCookies(t *testing.T) {
	s := NewGameState()
	s.addCookies(1)
	if s.GetCookies() != 1 {
		t.Errorf("Expected %e cookies, got %e", 1, s.GetCookies())
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

	b := newStandardBuilding(
		"New Building",
		"",
		nil,
		1,
	)
	s.loadBuildings(map[BuildingType]BuildingInterface{
		BUILDING_TYPE_CURSOR: b,
	})
	s.nBuildings[BUILDING_TYPE_CURSOR] = 1

	cps := s.calculateCPS()
	if cps != 1 {
		t.Errorf("Expected total CPS %e, got %e", 1, cps)
	}
}

func TestCalculateCPSSimpleUpgrade(t *testing.T) {
	s := NewGameState()

	b := newStandardBuilding(
		"New Building",
		"",
		nil,
		1,
	)

	s.loadBuildings(map[BuildingType]BuildingInterface{
		BUILDING_TYPE_CURSOR: b,
	})

	s.nBuildings[BUILDING_TYPE_CURSOR] = 1
	s.upgradeStatus[UPGRADE_ID_REINFORCED_INDEX_FINGER] = true

	u := newBuildingUpgrade(
		BUILDING_TYPE_CURSOR,
		"New Upgrade",
		"",
		0,
		2,
		0,
	)

	s.loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	cps := s.calculateCPS()
	if s.calculateCPS() != 2 {
		t.Errorf("Expected total CPS %e, got %e", 2, cps)
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

	u := newBuildingUpgrade(
		BUILDING_TYPE_CURSOR,
		"New Upgrade",
		"",
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

	u := newBuildingUpgrade(
		BUILDING_TYPE_CURSOR,
		"New Upgrade",
		"",
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

	u := newBuildingUpgrade(
		BUILDING_TYPE_CURSOR,
		"New Upgrade",
		"",
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

	b := newStandardBuilding(
		"New Building",
		"",
		nil,
		1,
	)

	s.loadBuildings(map[BuildingType]BuildingInterface{
		BUILDING_TYPE_CURSOR: b,
	})

	s.nBuildings[BUILDING_TYPE_CURSOR] = 1

	u := newBuildingUpgrade(
		BUILDING_TYPE_CURSOR,
		"New Upgrade",
		"",
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

	if s.GetCPS(start, end) != 2 {
		t.Errorf("Expected CPS %e, got %e", 2, s.GetCPS(start, end))
	}
}

func TestBuyCookiesPerClickUpgrade(t *testing.T) {
	s := NewGameState()

	s.loadCookiesPerClickRef(2)

	u := newBasicClickUpgrade(
		"New Upgrade",
		"",
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

	n := time.Now()
	if s.GetCookiesPerClick(n) != 6 {
		t.Errorf("Expected %e cookies per click, got %e", 6, s.GetCookiesPerClick(n))
	}
}

func TestBuyBuildingFree(t *testing.T) {
	s := NewGameState()

	b := newStandardBuilding(
		"New Building",
		"",
		func(current int) float64 { return 0 },
		1,
	)

	s.loadBuildings(map[BuildingType]BuildingInterface{
		BUILDING_TYPE_CURSOR: b,
	})

	s.BuyBuilding(BUILDING_TYPE_CURSOR)
	if s.GetCPS(start, end) != 1 {
		t.Errorf("Expected CPS %e, got %e", 1, s.GetCPS(start, end))
	}
}

func TestBuyNonexistentBuilding(t *testing.T) {
	s := NewGameState()

	if s.BuyBuilding(BUILDING_TYPE_CURSOR) {
		t.Error("Bought nonexistent building.")
	}
	if s.nBuildings[BUILDING_TYPE_CURSOR] != 0 {
		t.Error("Building was bought even though BuyBuilding returned not bought.")
	}
}

func TestBuyBuildingTooExpensive(t *testing.T) {
	s := NewGameState()

	b := newStandardBuilding(
		"New Building",
		"",
		func(current int) float64 { return 1 },
		1,
	)

	s.loadBuildings(map[BuildingType]BuildingInterface{
		BUILDING_TYPE_CURSOR: b,
	})

	if s.BuyBuilding(BUILDING_TYPE_CURSOR) {
		t.Error("Expected to not buy building, but bought anyways.")
	}
	if s.nBuildings[BUILDING_TYPE_CURSOR] != 0 {
		t.Errorf("Expected %d buildings, got %d", 0, s.nBuildings[BUILDING_TYPE_CURSOR])
	}
}

func TestBuyBuildingAffordable(t *testing.T) {
	s := NewGameState()
	s.addCookies(1)

	b := newStandardBuilding(
		"New Building",
		"",
		func(current int) float64 { return 1 },
		1,
	)

	s.loadBuildings(map[BuildingType]BuildingInterface{
		BUILDING_TYPE_CURSOR: b,
	})

	if !s.BuyBuilding(BUILDING_TYPE_CURSOR) {
		t.Error("Expected to buy building, but couldn't.")
	}
	if s.nBuildings[BUILDING_TYPE_CURSOR] != 1 {
		t.Errorf("Expected %d buildings, got %d", 1, s.nBuildings[BUILDING_TYPE_CURSOR])
	}
	if s.GetCookies() != 0 {
		t.Errorf("Expected %e cookies left, got %e", 0, s.GetCookies)
	}
}

func TestNewGameStateData(t *testing.T) {
	d := NewGameStateData()
	if d.Version != GAME_STATE_VERSION {
		t.Errorf("Expected %s version, got %s", GAME_STATE_VERSION, d.Version)
	}
}

func TestGameStateLoad(t *testing.T) {
	s := NewGameState()
	s.Load(GameStateData{
		Version: GAME_STATE_VERSION,
	})
}

func TestGameStateLoadBadVersion(t *testing.T) {
	s := NewGameState()
	err := s.Load(GameStateData{
		Version: "some-bad-version",
	})
	if err == nil {
		t.Error("Unexpected success while loading outdated data")
	}
}

func TestGameStateDump(t *testing.T) {
	s := NewGameState()
	s.Load(GameStateData{})

	s.addCookies(1)

	d := s.Dump()
	if d.NCookies != 1 {
		t.Errorf("Expected %e cookies, got %e", 1, d.NCookies)
	}
	if d.Version != GAME_STATE_VERSION {
		t.Errorf("Expected %s version, got %s", GAME_STATE_VERSION, d.Version)
	}
}
