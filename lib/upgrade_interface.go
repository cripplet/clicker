package cookie_clicker

type UpgradeInterface interface {
	GetName() string

	GetBuildingType() BuildingType
	GetBuildingMultiplier(g *GameStateStruct) float64
}
