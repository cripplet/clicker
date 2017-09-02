package cookie_clicker

type SimpleBuildingUpgrade struct {
	UpgradeInterface
	building_type       BuildingType
	name                string
	cost                float64
	building_multiplier float64
	minimum_buildings   int
}

func NewSimpleBuildingUpgrade(t BuildingType, n string, c float64, m float64, b int) *SimpleBuildingUpgrade {
	u := SimpleBuildingUpgrade{
		building_type:       t,
		name:                n,
		cost:                c,
		building_multiplier: m,
		minimum_buildings:   b,
	}
	return &u
}

func (self *SimpleBuildingUpgrade) GetIsUnlocked(g *GameStateStruct) bool {
	return (*g).GetNBuildings()[(*self).GetBuildingType()] >= (*self).minimum_buildings
}

func (self *SimpleBuildingUpgrade) GetCost(g *GameStateStruct) float64 {
	return (*self).cost
}

func (self *SimpleBuildingUpgrade) GetName() string {
	return (*self).name
}

func (self *SimpleBuildingUpgrade) GetBuildingType() BuildingType {
	return (*self).building_type
}

func (self *SimpleBuildingUpgrade) GetBuildingMultiplier(g *GameStateStruct) float64 {
	return (*self).building_multiplier
}
