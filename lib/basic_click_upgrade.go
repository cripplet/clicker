package cookie_clicker

type BasicClickUpgrade struct {
	UpgradeInterface
	name            string
	cost            float64
	clickMultiplier float64
}

func NewBasicClickUpgrade(n string, c float64, m float64) *BasicClickUpgrade {
	u := BasicClickUpgrade{
		name:            n,
		cost:            c,
		clickMultiplier: m,
	}
	return &u
}

func (self *BasicClickUpgrade) GetIsUnlocked(g *GameStateStruct) bool {
	return true
}

func (self *BasicClickUpgrade) GetCost(g *GameStateStruct) float64 {
	return self.cost
}

func (self *BasicClickUpgrade) GetName() string {
	return self.name
}

func (self *BasicClickUpgrade) GetBuildingType() BuildingType {
	return BUILDING_TYPE_ENUM_EOF
}

func (self *BasicClickUpgrade) GetBuildingMultiplier(g *GameStateStruct) float64 {
	return 1
}

func (self *BasicClickUpgrade) GetClickMultiplier(g *GameStateStruct) float64 {
	return self.clickMultiplier
}
