package cookie_clicker

type GameAPI interface {
	Load()
	Start()
	Stop()
	GetNBuildings() map[BuildingType]int
	GetUpgrades() map[UpgradeID]UpgradeInterface
	GetUpgradeStatus() map[UpgradeID]bool
	GetCookies() float64
	BuyUpgrade(UpgradeID) bool
	BuyBuilding(BuildingType) bool
	GetCPS() float64
	Click()
}
