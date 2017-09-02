package cookie_clicker

import (
	"time"
)

var EPOCH_MILLISECONDS time.Duration = time.Duration(time.Millisecond * 100)

type GameStateStruct struct {
	GameAPI
	n_cookies             float64
	cps                   float64 // cached value recalculated for each new upgrade
	n_buildings           map[BuildingType]int
	building_cps          map[BuildingType]float64
	building_cost         map[BuildingType]BuildingCostFunction
	upgrade_status        map[UpgradeID]bool
	upgrades              map[UpgradeID]UpgradeInterface
	main_loop_done_signal chan bool
}

func NewGameState() *GameStateStruct {
	g := GameStateStruct{
		n_buildings:    make(map[BuildingType]int),
		building_cps:   make(map[BuildingType]float64),
		building_cost:  make(map[BuildingType]BuildingCostFunction),
		upgrade_status: make(map[UpgradeID]bool),
		upgrades:       make(map[UpgradeID]UpgradeInterface),
	}

	var i BuildingType
	for i = 0; i < BUILDING_TYPE_ENUM_EOF; i++ {
		g.n_buildings[i] = 0
		g.building_cps[i] = 0
		g.building_cost[i] = func(current int) float64 { return 0 }
	}

	var j UpgradeID
	for j = 0; j < UPGRADE_ID_ENUM_EOF; j++ {
		g.upgrade_status[j] = false
	}

	return &g
}

/* Public API */

func (self *GameStateStruct) Load() {
	(*self).loadBuildingCost(BUILDING_COST_LOOKUP)
	(*self).loadUpgrades(UPGRADES_LOOKUP)
	(*self).loadBuildingCPS(BUILDING_CPS_LOOKUP)
}

func (self *GameStateStruct) Start() {
	(*self).main_loop_done_signal = make(chan bool)
	go (*self).startBlocking()
}

func (self *GameStateStruct) Stop() {
	(*self).main_loop_done_signal <- true
}

func (self *GameStateStruct) GetNBuildings() map[BuildingType]int {
	return (*self).n_buildings
}

func (self *GameStateStruct) GetUpgrades() map[UpgradeID]UpgradeInterface {
	return (*self).upgrades
}

func (self *GameStateStruct) GetUpgradeStatus() map[UpgradeID]bool {
	return (*self).upgrade_status
}

func (self *GameStateStruct) GetCookies() float64 {
	return (*self).n_cookies
}

func (self *GameStateStruct) BuyUpgrade(id UpgradeID) bool { // TODO(cripplet): Enforce upgrade cost check.
	_, present := (*self).upgrades[id]
	to_buy := present && !(*self).upgrade_status[id]
	if to_buy {
		(*self).upgrade_status[id] = true
		(*self).calculateCPS()
	}
	return to_buy
}

func (self *GameStateStruct) BuyBuilding(building_type BuildingType) bool {
	cost := (*self).building_cost[building_type]((*self).n_buildings[building_type])
	bought := (*self).subtractCookies(cost)
	if bought {
		(*self).n_buildings[building_type] += 1
		(*self).calculateCPS()
	}
	return bought
}

func (self *GameStateStruct) GetCPS() float64 {
	return (*self).cps
}

func (self *GameStateStruct) Click() { // TODO(cripplet): Add click upgrades.
	(*self).addCookies(1)
}

/* End public API */

func (self *GameStateStruct) setCPS(cps float64) {
	(*self).cps = cps
}

func (self *GameStateStruct) addCookies(n float64) {
	(*self).n_cookies += n
}

func (self *GameStateStruct) subtractCookies(n float64) bool {
	if (*self).n_cookies >= n {
		(*self).n_cookies -= n
		return true
	}
	return false
}

func (self *GameStateStruct) loadBuildingCost(c map[BuildingType]BuildingCostFunction) {
	for building_type := range (*self).building_cost {
		(*self).building_cost[building_type] = func(current int) float64 { return 0 }
	}

	for building_type, building_cost_function := range c {
		(*self).building_cost[building_type] = building_cost_function
	}
}

func (self *GameStateStruct) loadUpgrades(u map[UpgradeID]UpgradeInterface) {
	for upgrade_id := range (*self).upgrades {
		delete((*self).upgrades, upgrade_id)
	}

	for upgrade_id, upgrade_interface := range u {
		(*self).upgrades[upgrade_id] = upgrade_interface
	}
}

func (self *GameStateStruct) loadBuildingCPS(b map[BuildingType]float64) {
	for building_type := range (*self).building_cps {
		(*self).building_cps[building_type] = 0
	}

	for building_type, building_cps := range b {
		(*self).building_cps[building_type] = building_cps
	}
}

func (self *GameStateStruct) calculateCPS() float64 {
	building_cps_copy := make(map[BuildingType]float64)
	for building_type, building_type_cps := range (*self).building_cps {
		building_cps_copy[building_type] = building_type_cps
	}

	bought_upgrades := make([]UpgradeInterface, 0)
	for upgrade_id, upgrade := range (*self).upgrades {
		if (*self).upgrade_status[upgrade_id] {
			bought_upgrades = append(bought_upgrades, upgrade)
		}
	}

	for _, upgrade := range bought_upgrades {
		if upgrade.GetBuildingType() < BUILDING_TYPE_ENUM_EOF {
			building_cps_copy[upgrade.GetBuildingType()] *= upgrade.GetBuildingMultiplier(self)
		}
	}

	var total_cps float64
	for building_type, cps := range building_cps_copy {
		total_cps += float64((*self).n_buildings[building_type]) * cps
	}

	(*self).setCPS(total_cps)

	return total_cps
}

func calculateCookiesSince(start time.Time, end time.Time, cps float64) float64 {
	return cps * float64(end.Sub(start)) / float64(time.Second)
}

func (self *GameStateStruct) startBlocking() {
	t := time.NewTicker(EPOCH_MILLISECONDS)
	last_updated := time.Now()
	for {
		select {
		case <-t.C:
			current_time := time.Now()
			(*self).addCookies(calculateCookiesSince(last_updated, current_time, (*self).GetCPS()))
			last_updated = current_time
			continue
		case <-(*self).main_loop_done_signal:
			t.Stop()
			return
		}
	}
}
