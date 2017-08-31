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
	building_cps_scratch := make(map[BuildingType]float64)
	for k, v := range building_cps {
		building_cps_scratch[k] = v
	}

	for upgrade_id, upgrade := range u {
		if is_bought, ok := (*self).upgrade_status[upgrade_id]; ok {
			if is_bought {
				if upgrade.GetBuildingType() < BUILDING_TYPE_ENUM_EOF {
					building_cps_scratch[upgrade.GetBuildingType()] *= upgrade.GetBuildingMultiplier(self)
				}
			}
		}
	}

	var cps float64
	for building_type, building_cps := range building_cps_scratch {
		cps += float64((*self).n_buildings[building_type]) * building_cps
	}

	(*self).cps = cps

	return cps
}
