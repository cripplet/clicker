package cookie_clicker

type SimpleBuildingUpgrade struct {
	UpgradeInterface
	building_type       BuildingType
	name                string
	building_multiplier float64
}

func NewSimpleBuildingUpgrade(t BuildingType, n string, m float64) *SimpleBuildingUpgrade {
	u := SimpleBuildingUpgrade{
		building_type:       t,
		name:                n,
		building_multiplier: m,
	}
	return &u
}

func (self *SimpleBuildingUpgrade) GetName() string {
	return self.name
}

func (self *SimpleBuildingUpgrade) GetBuildingType() BuildingType {
	return (*self).building_type
}

func (self *SimpleBuildingUpgrade) GetBuildingMultiplier(g *GameStateStruct) float64 {
	return (*self).building_multiplier
}
