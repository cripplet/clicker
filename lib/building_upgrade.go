package cookie_clicker

type BuildingUpgrade struct {
	UpgradeInterface
	building_type       BuildingType
	name                string
	cost                float64
	building_multiplier float64
	minimum_buildings   int
}

func NewBuildingUpgrade(t BuildingType, n string, c float64, m float64, b int) *BuildingUpgrade {
	u := BuildingUpgrade{
		building_type:       t,
		name:                n,
		cost:                c,
		building_multiplier: m,
		minimum_buildings:   b,
	}
	return &u
}

func (self *BuildingUpgrade) GetIsUnlocked(g *GameStateStruct) bool {
	return (*g).GetNBuildings()[(*self).GetBuildingType()] >= (*self).minimum_buildings
}

func (self *BuildingUpgrade) GetCost(g *GameStateStruct) float64 {
	return (*self).cost
}

func (self *BuildingUpgrade) GetName() string {
	return (*self).name
}

func (self *BuildingUpgrade) GetBuildingType() BuildingType {
	return (*self).building_type
}

func (self *BuildingUpgrade) GetBuildingMultiplier(g *GameStateStruct) float64 {
	return (*self).building_multiplier
}

func (self *BuildingUpgrade) GetClickMultiplier(g *GameStateStruct) float64 {
	return 1
}
