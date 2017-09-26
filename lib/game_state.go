package cookie_clicker

import (
	"errors"
	"fmt"
	"time"
)

var GAME_STATE_VERSION string = "v0.01"

type GameStateData struct {
	Version       string
	NCookies      float64
	NBuildings    map[BuildingType]int
	UpgradeStatus map[UpgradeID]bool
}

type GameStateStruct struct {
	/* Imported Fields */
	nCookies      float64
	nBuildings    map[BuildingType]int
	upgradeStatus map[UpgradeID]bool
	/* Calculated Cache Fields */
	cookiesPerClick float64
	cps             float64
	/* Immutable Fields */
	cookiesPerClickRef float64
	upgrades           map[UpgradeID]UpgradeInterface
	buildings          map[BuildingType]BuildingInterface
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

/* Public API */

func (self *GameStateStruct) Load(d GameStateData) error {
	err := self.loadData(d)
	if err != nil {
		return err
	}
	self.loadBuildings(BUILDINGS_LOOKUP)
	self.loadUpgrades(UPGRADES_LOOKUP)
	self.loadCookiesPerClickRef(COOKIES_PER_CLICK_LOOKUP)

	self.setCookiesPerClick(self.calculateCookiesPerClick())
	self.setCPS(self.calculateCPS())

	return nil
}

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

func (self *GameStateStruct) GetNBuildings() map[BuildingType]int {
	return self.nBuildings
}

func (self *GameStateStruct) GetBuildings() map[BuildingType]BuildingInterface {
	return self.buildings
}

func (self *GameStateStruct) GetUpgrades() map[UpgradeID]UpgradeInterface {
	return self.upgrades
}

func (self *GameStateStruct) GetUpgradeStatus() map[UpgradeID]bool {
	return self.upgradeStatus
}

func (self *GameStateStruct) GetCookies() float64 {
	return self.nCookies
}

func (self *GameStateStruct) BuyUpgrade(id UpgradeID) bool { // TODO(cripplet): Enforce upgrade cost check.
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

func (self *GameStateStruct) GetCPS(start time.Time, end time.Time) float64 { // TODO(cripplet): Calculate timed buffs here.
	return self.cps * float64(end.Sub(start)) / float64(time.Second)
}

func (self *GameStateStruct) GetCookiesPerClick() float64 { // TODO(cripplet): Calculate timed buffs here.
	return self.cookiesPerClick
}

func (self *GameStateStruct) MineCookies(start time.Time, end time.Time) {
	self.addCookies(self.GetCPS(start, end))
}

func (self *GameStateStruct) Click() {
	self.addCookies(self.GetCookiesPerClick())
}

/* End public API */

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
		if upgrade.GetBuildingType() < BUILDING_TYPE_ENUM_EOF {
			buildingCPSRefCopy[upgrade.GetBuildingType()] *= upgrade.GetBuildingMultiplier(self)
		}
	}

	var totalCPS float64
	for buildingType, cps := range buildingCPSRefCopy {
		totalCPS += float64(self.nBuildings[buildingType]) * cps
	}

	return totalCPS
}
