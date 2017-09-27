package cookie_clicker

type basicClickUpgrade struct {
	UpgradeInterface
	name            string
	cost            float64
	clickMultiplier float64
}

func newBasicClickUpgrade(n string, c float64, m float64) *basicClickUpgrade {
	u := basicClickUpgrade{
		name:            n,
		cost:            c,
		clickMultiplier: m,
	}
	return &u
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

func (self *basicClickUpgrade) GetBuildingMultiplier(g *GameStateStruct) float64 {
	return 1
}

func (self *basicClickUpgrade) GetClickMultiplier(g *GameStateStruct) float64 {
	return self.clickMultiplier
}
