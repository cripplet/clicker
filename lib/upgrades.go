package cookie_clicker

var UPGRADES_LOOKUP map[UpgradeID]UpgradeInterface = map[UpgradeID]UpgradeInterface{
	UPGRADE_ID_REINFORCED_INDEX_FINGER: NewBasicClickUpgrade(
		"Reinforced Index Finger",
		100,
		2,
	),
}
