package cookie_clicker

type UpgradeInterface interface {
	GetName() string
	GetCost(g *GameStateStruct) float64
	GetIsUnlocked(g *GameStateStruct) bool

	GetBuildingType() BuildingType
	GetBuildingMultiplier(g *GameStateStruct) float64

	GetClickMultiplier(g *GameStateStruct) float64
}
