package cookie_clicker

type buildingUpgrade struct {
	UpgradeInterface
	buildingType       BuildingType
	name               string
	description        string
	cost               float64
	buildingMultiplier float64
	minimumBuildings   int
}

func newBuildingUpgrade(t BuildingType, n string, d string, c float64, m float64, b int) *buildingUpgrade {
	u := buildingUpgrade{
		buildingType:       t,
		name:               n,
		description:        d,
		cost:               c,
		buildingMultiplier: m,
		minimumBuildings:   b,
	}
	return &u
}

func (self *buildingUpgrade) GetDescription() string {
	return self.description
}

func (self *buildingUpgrade) GetIsUnlocked(g *GameStateStruct) bool {
	return (*g).GetNBuildings()[self.buildingType] >= self.minimumBuildings
}

func (self *buildingUpgrade) GetCost(g *GameStateStruct) float64 {
	return self.cost
}

func (self *buildingUpgrade) GetName() string {
	return self.name
}

func (self *buildingUpgrade) GetBuildingMultipliers(g *GameStateStruct) map[BuildingType]float64 {
	c := map[BuildingType]float64{}
	var i BuildingType
	for i = BuildingType(0); i < BUILDING_TYPE_ENUM_EOF; i++ {
		if self.buildingType == i {
			c[i] = self.buildingMultiplier
		} else {
			c[i] = 1
		}
	}
	return c
}

func (self *buildingUpgrade) GetClickMultiplier(g *GameStateStruct) float64 {
	return 1
}
