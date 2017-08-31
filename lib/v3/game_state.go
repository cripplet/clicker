package cookie_clicker

type BuildingType int
type UpgradeID int

const (
	BUILDING_TYPE_MOUSE BuildingType = iota
	BUILDING_TYPE_GRANDMA

	BUILDING_TYPE_ENUM_EOF
)

const (
	UPGRADE_ID_REINFORCED_INDEX_FINGER UpgradeID = iota
	UPGRADE_ID_THOUSAND_FINGERS
	UPGRADE_ID_FORWARDS_FROM_GRANDMA
	// UPGRADE_ID_PLAIN_COOKIES

	UPGRADE_ID_ENUM_EOF
)

type UpgradeInterface interface {
	GetName() string
	GetCost(g *GameStateStruct) float64
	GetUnlockStatus(g *GameStateStruct) bool

	GetMultiplicativeCPSContribution(g *GameStateStruct) float64
	GetAdditiveCPSContribution(g *GameStateStruct) float64

	GetBuildingType() BuildingType
	GetBuildingMultiplier(g *GameStateStruct) float64
}

var building_cps map[BuildingType]float64 = map[BuildingType]float64{
	BUILDING_TYPE_MOUSE: 0.2,
	BUILDING_TYPE_GRANDMA: 1,
}

type GameStateStruct struct {
	n_cookies float64
	cps float64  // cached value recalculated for each new upgrade
	n_buildings map[BuildingType]int
	upgrade_status map[UpgradeID]bool
}

func NewGameState() *GameStateStruct {
	g := GameStateStruct{
		n_buildings: make(map[BuildingType]int),
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

	for upgrade_id, upgrade := range u {
		if is_bought, ok := (*self).upgrade_status[upgrade_id]; ok {
			if is_bought {
				cps += upgrade.GetAdditiveCPSContribution(self)
			}
		}
	}

	for upgrade_id, upgrade := range u {
		if is_bought, ok := (*self).upgrade_status[upgrade_id]; ok {
			if is_bought {
				cps *= upgrade.GetMultiplicativeCPSContribution(self)
			}
		}
	}

	(*self).cps = cps

	return cps
}

/* SimpleUpgrade */

type SimpleBuildingUpgrade struct {
	UpgradeInterface
	building_type BuildingType
	name string
	cost float64
	building_multiplier float64
	minimum_buildings int
}

func NewSimpleBuildingUpgrade(t BuildingType, n string, c float64, b float64, m int) *SimpleBuildingUpgrade {
	u := SimpleBuildingUpgrade{
		building_type: t,
		name: n,
		cost: c,
		building_multiplier: b,
		minimum_buildings: m,
	}
	return &u
}

func (self *SimpleBuildingUpgrade) GetName() string {
	return self.name
}

func (self *SimpleBuildingUpgrade) GetCost(g *GameStateStruct) float64 {
	return self.cost
}

func (self *SimpleBuildingUpgrade) GetUnlockStatus(g *GameStateStruct) bool {
	// TODO(cripplet): Lock
	return g.n_buildings[self.building_type] >= self.minimum_buildings
}

func (self *SimpleBuildingUpgrade) GetMultiplicativeCPSContribution(g *GameStateStruct) float64 {
	return 1
}

func (self *SimpleBuildingUpgrade) GetAdditiveCPSContribution(g *GameStateStruct) float64 {
	return 0
}

func (self *SimpleBuildingUpgrade) GetBuildingType() BuildingType {
	return (*self).building_type
}

func (self *SimpleBuildingUpgrade) GetBuildingMultiplier(g *GameStateStruct) float64 {
	return 2
}

/* END SimpleUpgrade */
