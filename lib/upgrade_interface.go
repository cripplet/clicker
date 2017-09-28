package cookie_clicker

// UpgradeInterface is the public API for all upgrades.
//
// TODO(cripplet): Add GetCPSMultiplier.
type UpgradeInterface interface {
	GetName() string
	GetCost(g *GameStateStruct) float64
	GetDescription() string

	// Returns true if the player can buy the upgrade.
	GetIsUnlocked(g *GameStateStruct) bool

	// Returns the CPS multiplier of the given building type.
	GetBuildingMultipliers(g *GameStateStruct) map[BuildingType]float64

	// Return the multiplier for cookies added per player click.
	GetClickMultiplier(g *GameStateStruct) float64
}
