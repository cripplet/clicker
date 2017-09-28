// Package cookie_clicker contains all data structs and functions used to
// run a game of CookieClicker.
package cookie_clicker

import (
	"errors"
	"fmt"
	"time"
)

var GAME_STATE_VERSION string = "v0.01"

// GameStateData is the library representation of the data necessary to import
// and export a game.
type GameStateData struct {
	Version       string               `json:"version"`
	NCookies      float64              `json:"n_cookies"`
	NBuildings    map[BuildingType]int `json:"buildings"`
	UpgradeStatus map[UpgradeID]bool   `json:"upgrades"`
}

// GameStateStruct represents a game of CookieClicker.
type GameStateStruct struct {
	// Number of cookies the player currently owns. This imported from
	// GameStateData.
	nCookies float64

	// Number of buildings the player currently owns of each type. This is
	// imported from GameStateData.
	nBuildings map[BuildingType]int

	// Upgrades that the player currently owns. This is imported from
	// GameStateData.
	upgradeStatus map[UpgradeID]bool

	// Cookies added to the player's bank per click.
	cookiesPerClick float64

	// Cookies added to the player's bank per second.
	cps float64

	// Base cookies added per click. This is used to calculate the actual
	// cookiesPerClick (taking upgrades, buildings, etc. into effect).
	cookiesPerClickRef float64

	// Copy of all upgrades the game knows about.
	upgrades map[UpgradeID]UpgradeInterface

	// Copy of all building types the game knows about.
	buildings map[BuildingType]BuildingInterface
}

func NewGameStateData() *GameStateData {
	d := GameStateData{
		Version:       GAME_STATE_VERSION,
		NBuildings:    make(map[BuildingType]int),
		UpgradeStatus: make(map[UpgradeID]bool),
	}
	return &d
}

func NewGameState() *GameStateStruct {
	g := GameStateStruct{
		nBuildings:    make(map[BuildingType]int),
		upgradeStatus: make(map[UpgradeID]bool),
		upgrades:      make(map[UpgradeID]UpgradeInterface),
		buildings:     make(map[BuildingType]BuildingInterface),
	}

	var i BuildingType
	for i = 0; i < BUILDING_TYPE_ENUM_EOF; i++ {
		g.nBuildings[i] = 0
	}

	var j UpgradeID
	for j = 0; j < UPGRADE_ID_ENUM_EOF; j++ {
		g.upgradeStatus[j] = false
	}

	return &g
}

// Load a game given import data. This will also populate other hidden
// reference fields.
func (self *GameStateStruct) Load(d GameStateData) error {
	self.loadBuildings(BUILDINGS_LOOKUP)
	self.loadUpgrades(UPGRADES_LOOKUP)

	err := self.loadData(d)
	if err != nil {
		return err
	}

	self.loadCookiesPerClickRef(COOKIES_PER_CLICK_LOOKUP)

	self.setCookiesPerClick(self.calculateCookiesPerClick())
	self.setCPS(self.calculateCPS())

	return nil
}

// Export game data.
func (self *GameStateStruct) Dump() GameStateData {
	d := GameStateData{
		Version:       GAME_STATE_VERSION,
		NCookies:      self.nCookies,
		NBuildings:    make(map[BuildingType]int),
		UpgradeStatus: make(map[UpgradeID]bool),
	}
	for buildingType, nBuildings := range self.nBuildings {
		d.NBuildings[buildingType] = nBuildings
	}
	for upgradeID, bought := range self.upgradeStatus {
		d.UpgradeStatus[upgradeID] = bought
	}
	return d
}

// Get the current number of buildings per type that the player currently owns.
func (self *GameStateStruct) GetNBuildings() map[BuildingType]int {
	return self.nBuildings
}

// Get a list of buildings the game knows about.
func (self *GameStateStruct) GetBuildings() map[BuildingType]BuildingInterface {
	return self.buildings
}

// Get a list of ugrades the game knows about.
func (self *GameStateStruct) GetUpgrades() map[UpgradeID]UpgradeInterface {
	return self.upgrades
}

// Get a list of upgrades and whether or not the player currently owns it.
func (self *GameStateStruct) GetUpgradeStatus() map[UpgradeID]bool {
	return self.upgradeStatus
}

// Get the current number of cookies the player owns.
func (self *GameStateStruct) GetCookies() float64 {
	return self.nCookies
}

// Attempts to buy an upgrade and subtract the cost of the upgrade from the
// player bank. BuyUpgrade will return true if this was successful. This will
// return false if the player has insufficient funds, if the upgrade is
// already bought, or if the upgrade is not unlocked.
func (self *GameStateStruct) BuyUpgrade(id UpgradeID) bool {
	upgrade, present := self.upgrades[id]
	to_buy := present && !self.upgradeStatus[id]
	bought := to_buy && upgrade.GetIsUnlocked(self) && self.subtractCookies(upgrade.GetCost(self))
	if bought {
		self.upgradeStatus[id] = true
		self.setCPS(self.calculateCPS())
		self.setCookiesPerClick(self.calculateCookiesPerClick())

	}
	return bought
}

// Attempts to buy a building and subtract the cost of the building from the
// player bank. BuyBuilding will return true if this was successful. This will
// r eturn false if the player has insufficient funds.
func (self *GameStateStruct) BuyBuilding(buildingType BuildingType) bool {
	building, present := self.buildings[buildingType]
	bought := present && self.subtractCookies(building.GetCost(self.nBuildings[buildingType]+1))
	if bought {
		self.nBuildings[buildingType] += 1
		self.setCPS(self.calculateCPS())
		self.setCookiesPerClick(self.calculateCookiesPerClick())
	}
	return bought
}

// Get the number of cookies added to the player's bank over a given time
// period.
//
// TODO(cripplet): Calculate timed buffs here by adding a list of time events.
func (self *GameStateStruct) GetCPS(start time.Time, end time.Time) float64 {
	return self.cps * float64(end.Sub(start)) / float64(time.Second)
}

// Get the number of cookies added to the player's bank per click.
//
// TODO(cripplet): Calculate timed buffs here by adding a list of time events
// and add a time parameter to the function signature.
func (self *GameStateStruct) GetCookiesPerClick(clickTime time.Time) float64 {
	return self.cookiesPerClick
}

// Commit cookies from CPS contribution to the bank.
func (self *GameStateStruct) MineCookies(start time.Time, end time.Time) {
	self.addCookies(self.GetCPS(start, end))
}

// Commit cookies from physical click contributions to the bank.
func (self *GameStateStruct) Click(clickTime time.Time) {
	self.addCookies(self.GetCookiesPerClick(clickTime))
}

func (self *GameStateStruct) setCPS(cps float64) {
	self.cps = cps
}

func (self *GameStateStruct) addCookies(n float64) {
	self.nCookies += n
}

func (self *GameStateStruct) subtractCookies(n float64) bool {
	if self.nCookies >= n {
		self.nCookies -= n
		return true
	}
	return false
}

func (self *GameStateStruct) setCookiesPerClick(c float64) {
	self.cookiesPerClick = c
}

func (self *GameStateStruct) loadData(d GameStateData) error {
	if d.Version != GAME_STATE_VERSION {
		return errors.New(fmt.Sprintf("Outdated data version: %s < %s", d.Version, GAME_STATE_VERSION))
	}
	self.nCookies = d.NCookies

	for buildingType, _ := range self.nBuildings {
		self.nBuildings[buildingType] = d.NBuildings[buildingType]
	}
	for upgradeType, _ := range self.upgradeStatus {
		self.upgradeStatus[upgradeType] = d.UpgradeStatus[upgradeType]
	}
	return nil
}

func (self *GameStateStruct) loadBuildings(b map[BuildingType]BuildingInterface) {
	for buildingType := range self.buildings {
		delete(self.buildings, buildingType)
	}

	for buildingType, buildingInterface := range b {
		self.buildings[buildingType] = buildingInterface
	}
}

func (self *GameStateStruct) loadUpgrades(u map[UpgradeID]UpgradeInterface) {
	for upgradeID := range self.upgrades {
		delete(self.upgrades, upgradeID)
	}

	for upgradeID, upgradeInterface := range u {
		self.upgrades[upgradeID] = upgradeInterface
	}
}

func (self *GameStateStruct) loadCookiesPerClickRef(c float64) {
	self.cookiesPerClickRef = c
}

func (self *GameStateStruct) calculateCookiesPerClick() float64 {
	cookiesPerClickCopy := self.cookiesPerClickRef
	for upgradeID, bought := range self.upgradeStatus {
		if bought {
			cookiesPerClickCopy *= self.GetUpgrades()[upgradeID].GetClickMultiplier(self)
		}
	}

	return cookiesPerClickCopy
}

func (self *GameStateStruct) calculateCPS() float64 {
	buildingCPSRefCopy := make(map[BuildingType]float64)
	for buildingType, building := range self.buildings {
		buildingCPSRefCopy[buildingType] = building.GetCPS()
	}

	boughtUpgrades := make([]UpgradeInterface, 0)
	for upgradeID, upgrade := range self.upgrades {
		if self.upgradeStatus[upgradeID] {
			boughtUpgrades = append(boughtUpgrades, upgrade)
		}
	}

	for _, upgrade := range boughtUpgrades {
		for buildingType, buildingMultiplier := range upgrade.GetBuildingMultipliers(self) {
			buildingCPSRefCopy[buildingType] *= buildingMultiplier
		}
	}

	var totalCPS float64
	for buildingType, cps := range buildingCPSRefCopy {
		totalCPS += float64(self.nBuildings[buildingType]) * cps
	}

	return totalCPS
}
