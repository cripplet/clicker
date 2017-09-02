package cookie_clicker

import (
	"fmt"
	"testing"
	"time"
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

func TestAddCookies(t *testing.T) {
	s := NewGameState()
	(*s).addCookies(1)
	if (*s).GetCookies() != 1 {
		t.Error(fmt.Sprintf("Expected %e cookies, got %e", 1, (*s).GetCookies()))
	}
}

func TestSubtractCookiesTooExpensive(t *testing.T) {
	s := NewGameState()
	if (*s).subtractCookies(1) {
		t.Error("Subtracted cookies when not possible")
	}
}

func TestSubtractCookies(t *testing.T) {
	s := NewGameState()
	(*s).addCookies(1)
	if !(*s).subtractCookies(1) {
		t.Error("Could not subtract cookies")
	}
	if (*s).GetCookies() != 0 {
		t.Error("Expected %e cookies, got %e", 0, (*s).GetCookies())
	}
}

func TestCalculateCPSNoUpgrades(t *testing.T) {
	s := NewGameState()

	(*s).loadBuildingCPS(BUILDING_CPS_LOOKUP)
	(*s).n_buildings[BUILDING_TYPE_MOUSE] = 1

	if (*s).calculateCPS() != 0.2 {
		t.Error(fmt.Sprintf("Expected total CPS %e, got %e", 0.2, (*s).GetCPS()))
	}
}

func TestCalculateCPSSimpleUpgrade(t *testing.T) {
	s := NewGameState()

	(*s).loadBuildingCPS(BUILDING_CPS_LOOKUP)
	(*s).n_buildings[BUILDING_TYPE_MOUSE] = 1
	(*s).upgrade_status[UPGRADE_ID_REINFORCED_INDEX_FINGER] = true

	u := NewSimpleBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		0,
		2,
		0,
	)

	(*s).loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if (*s).calculateCPS() != 0.4 {
		t.Error(fmt.Sprintf("Expected total CPS %e, got %e", 0.4, (*s).GetCPS()))
	}
}

func TestBuyNonexistentUpgrade(t *testing.T) {
	s := NewGameState()

	if (*s).BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Bought nonexistent upgrade.")
	}
	if (*s).upgrade_status[UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("Upgrade was bought even though BuyUpgrade returned not bought.")
	}
}

func TestDoubleBuyUpgrade(t *testing.T) {
	s := NewGameState()

	u := NewSimpleBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		100,
		2,
		0,
	)

	(*s).loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	(*s).upgrade_status[UPGRADE_ID_REINFORCED_INDEX_FINGER] = true

	if (*s).BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Double bought upgrade.")
	}
	if !(*s).upgrade_status[UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("Unbought upgrade.")
	}
}

func TestBuyUpgradeTooExpensive(t *testing.T) {
	s := NewGameState()

	u := NewSimpleBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		100,
		2,
		0,
	)

	(*s).loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if (*s).BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Bought too expensive an upgrade.")
	}
}

func TestBuyUpgradeLocked(t *testing.T) {
	s := NewGameState()

	u := NewSimpleBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		0,
		2,
		1,
	)

	(*s).loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if (*s).BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Bought a locked upgrade")
	}
}

func TestBuyUpgrade(t *testing.T) {
	s := NewGameState()

	(*s).loadBuildingCPS(BUILDING_CPS_LOOKUP)
	(*s).n_buildings[BUILDING_TYPE_MOUSE] = 1

	u := NewSimpleBuildingUpgrade(
		BUILDING_TYPE_MOUSE,
		"New Upgrade",
		0,
		2,
		0,
	)

	(*s).loadUpgrades(map[UpgradeID]UpgradeInterface{
		UPGRADE_ID_REINFORCED_INDEX_FINGER: u,
	})

	if !(*s).BuyUpgrade(UPGRADE_ID_REINFORCED_INDEX_FINGER) {
		t.Error("Could not buy simple upgrade.")
	}

	if !(*s).upgrade_status[UPGRADE_ID_REINFORCED_INDEX_FINGER] {
		t.Error("BuyUpgrade reported upgrade was bought, but does not show up as bought.")
	}

	if (*s).GetCPS() != 0.4 {
		t.Error(fmt.Sprintf("Expected CPS %e, got %e", 0.4, (*s).GetCPS()))
	}
}

func TestBuyBuildingFree(t *testing.T) {
	s := NewGameState()
	(*s).loadBuildingCPS(BUILDING_CPS_LOOKUP)
	(*s).BuyBuilding(BUILDING_TYPE_MOUSE)
	if (*s).GetCPS() != 0.2 {
		t.Error(fmt.Sprintf("Expected CPS %e, got %e", 0.2, (*s).GetCPS()))
	}
}

func TestBuyBuildingTooExpensive(t *testing.T) {
	s := NewGameState()
	(*s).loadBuildingCost(map[BuildingType]BuildingCostFunction{
		BUILDING_TYPE_MOUSE: func(current int) float64 { return 1 },
	})
	if (*s).BuyBuilding(BUILDING_TYPE_MOUSE) {
		t.Error("Expected to not buy building, but bought anyways.")
	}
	if (*s).n_buildings[BUILDING_TYPE_MOUSE] != 0 {
		t.Error(fmt.Sprintf("Expected %d buildings, got %d", 0, (*s).n_buildings[BUILDING_TYPE_MOUSE]))
	}
}

func TestBuyBuildingAffordable(t *testing.T) {
	s := NewGameState()
	(*s).addCookies(1)
	(*s).loadBuildingCost(map[BuildingType]BuildingCostFunction{
		BUILDING_TYPE_MOUSE: func(current int) float64 { return 1 },
	})
	if !(*s).BuyBuilding(BUILDING_TYPE_MOUSE) {
		t.Error("Expected to buy building, but couldn't.")
	}
	if (*s).n_buildings[BUILDING_TYPE_MOUSE] != 1 {
		t.Error(fmt.Sprintf("Expected %d buildings, got %d", 1, (*s).n_buildings[BUILDING_TYPE_MOUSE]))
	}
	if (*s).GetCookies() != 0 {
		t.Error(fmt.Sprintf("Expected %e cookies left, got %e", 0, (*s).GetCookies))
	}
}

func TestCalculateCookiesSinceNoCPS(t *testing.T) {
	c := calculateCookiesSince(time.Now(), time.Now().Add(time.Second), 0)
	if c != 0 {
		t.Error(fmt.Sprintf("Expected %e cookies with no CPS, got %e", 0, c))
	}
}

func TestCalculateCookiesSinceNoDuration(t *testing.T) {
	n := time.Now()
	c := calculateCookiesSince(n, n, 1)
	if c != 0 {
		t.Error(fmt.Sprintf("Expected %e cookies with no time passing, got %e", 0, c))
	}
}

func TestCalculateCookiesSince(t *testing.T) {
	n := time.Now()
	c := calculateCookiesSince(n, n.Add(time.Second), 1)
	if c != 1 {
		t.Error(fmt.Sprintf("Expected %e cookies with 1 CPS, got %e", 1, c))
	}
}

func TestGameStateStartStop(t *testing.T) {
	s := NewGameState()
	(*s).setCPS(1)
	(*s).Start()
	time.Sleep(2 * time.Second)
	(*s).Stop()
	if (*s).GetCookies() < (*s).GetCPS() {
		t.Error(fmt.Sprintf("Expected %e cookies, got %e", (*s).GetCPS(), (*s).GetCookies()))
	}
}

func TestGameStateLoad(t *testing.T) {
	s := NewGameState()
	s.Load()

	if (*s).building_cps[BUILDING_TYPE_MOUSE] != BUILDING_CPS_LOOKUP[BUILDING_TYPE_MOUSE] {
		t.Error("Could not load game state constants properly")
	}
}
