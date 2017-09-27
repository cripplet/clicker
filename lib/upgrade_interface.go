package cookie_clicker

// UpgradeInterface is the public API for all upgrades.
//
// TOOD(cripplet): Add GetDescription.
//
// TODO(cripplet): Add GetCPSMultiplier.
//
// TODO(cripplet): Change GetBuildingMultiplier -> GetBuildingMultipliers(...) map[BuildingType]float64
//
// TODO(cripplet): Remove GetBuildingType
type UpgradeInterface interface {
	GetName() string
	GetCost(g *GameStateStruct) float64

	// Returns true if the player can buy the upgrade.
	GetIsUnlocked(g *GameStateStruct) bool

	// Returns the type of the building this upgrade is associated with.
	// It is acceptable for the return value to be BUILDING_TYPE_ENUM_EOF
	GetBuildingType() BuildingType

	// Returns the CPS multiplier of the given building type.
	GetBuildingMultiplier(g *GameStateStruct) float64

	// Return the multiplier for cookies added per player click.
	GetClickMultiplier(g *GameStateStruct) float64
}
