package cookie_clicker

type GameStateStruct struct {
	n_cookies      float64
	cps            float64 // cached value recalculated for each new upgrade
	n_buildings    map[BuildingType]int
	upgrade_status map[UpgradeID]bool
}

func NewGameState() *GameStateStruct {
	g := GameStateStruct{
		n_buildings:    make(map[BuildingType]int),
		upgrade_status: make(map[UpgradeID]bool),
	}

	var i BuildingType
	for i = 0; i < BUILDING_TYPE_ENUM_EOF; i++ {
		g.n_buildings[i] = 0
	}

	var j UpgradeID
	for j = 0; j < UPGRADE_ID_ENUM_EOF; j++ {
		g.upgrade_status[j] = false
	}

	return &g
}

func (self *GameStateStruct) CalculateCPS(u map[UpgradeID]UpgradeInterface) float64 {
	building_cps_copy := make(map[BuildingType]float64)
	for building_type, building_type_cps := range BUILDING_CPS_LOOKUP {
		building_cps_copy[building_type] = building_type_cps
	}

	bought_upgrades := make([]UpgradeInterface, 0)
	for upgrade_id, upgrade := range u {
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

	(*self).cps = total_cps

	return total_cps
}
