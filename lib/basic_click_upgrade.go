package cookie_clicker

type basicClickUpgrade struct {
	UpgradeInterface
	name            string
	description     string
	cost            float64
	clickMultiplier float64
}

func newBasicClickUpgrade(n string, d string, c float64, m float64) *basicClickUpgrade {
	u := basicClickUpgrade{
		name:            n,
		description:     d,
		cost:            c,
		clickMultiplier: m,
	}
	return &u
}

func (self *basicClickUpgrade) GetDescription() string {
	return self.description
}

func (self *basicClickUpgrade) GetIsUnlocked(g *GameStateStruct) bool {
	return true
}

func (self *basicClickUpgrade) GetCost(g *GameStateStruct) float64 {
	return self.cost
}

func (self *basicClickUpgrade) GetName() string {
	return self.name
}

func (self *basicClickUpgrade) GetBuildingType() BuildingType {
	return BUILDING_TYPE_ENUM_EOF
}

func (self *basicClickUpgrade) GetBuildingMultipliers(g *GameStateStruct) map[BuildingType]float64 {
	c := map[BuildingType]float64{}
	var i BuildingType
	for i = BuildingType(0); i < BUILDING_TYPE_ENUM_EOF; i++ {
		c[i] = 1
	}
	return c
}

func (self *basicClickUpgrade) GetClickMultiplier(g *GameStateStruct) float64 {
	return self.clickMultiplier
}
