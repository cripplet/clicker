package cookie_clicker

type BuildingUpgrade struct {
	UpgradeInterface
	buildingType       BuildingType
	name               string
	cost               float64
	buildingMultiplier float64
	minimumBuildings   int
}

func NewBuildingUpgrade(t BuildingType, n string, c float64, m float64, b int) *BuildingUpgrade {
	u := BuildingUpgrade{
		buildingType:       t,
		name:               n,
		cost:               c,
		buildingMultiplier: m,
		minimumBuildings:   b,
	}
	return &u
}

func (self *BuildingUpgrade) GetIsUnlocked(g *GameStateStruct) bool {
	return (*g).GetNBuildings()[self.GetBuildingType()] >= self.minimumBuildings
}

func (self *BuildingUpgrade) GetCost(g *GameStateStruct) float64 {
	return self.cost
}

func (self *BuildingUpgrade) GetName() string {
	return self.name
}

func (self *BuildingUpgrade) GetBuildingType() BuildingType {
	return self.buildingType
}

func (self *BuildingUpgrade) GetBuildingMultiplier(g *GameStateStruct) float64 {
	return self.buildingMultiplier
}

func (self *BuildingUpgrade) GetClickMultiplier(g *GameStateStruct) float64 {
	return 1
}
