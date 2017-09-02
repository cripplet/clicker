package cookie_clicker

type GameAPI interface {
	// Load pre-defined game state constants.
	Load()

	// Start the main game loop. This is non-blocking.
	Start()

	// Stop the main game loop. Blocks until loop terminates.
	Stop()

	// Get current buildings.
	GetNBuildings() map[BuildingType]int

	// Get list of upgrades.
	GetUpgrades() map[UpgradeID]UpgradeInterface

	// Get upgrade purchase status.
	GetUpgradeStatus() map[UpgradeID]bool

	// Get building cost functions.
	GetBuildingCost() map[BuildingType]BuildingCostFunction

	// Get number of cookies.
	GetCookies() float64

	// Attempt to buy an upgrade.
	BuyUpgrade(UpgradeID) bool

	// Attempt to buy one building.
	BuyBuilding(BuildingType) bool

	// Get current overall CPS.
	GetCPS() float64

	// Get cookies per click.
	GetCookiesPerClick() float64

	// Click a giant cookie.
	Click() float64
}
