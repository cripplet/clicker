package cookie_clicker

var UPGRADES_LOOKUP map[UpgradeID]UpgradeInterface = map[UpgradeID]UpgradeInterface{
	UPGRADE_ID_REINFORCED_INDEX_FINGER: newBasicClickUpgrade(
		"Reinforced Index Finger",
		// TODO(cripplet): Actually implement this.
		"The mouse and cursors are twice as efficient.",
		100,
		2,
	),
}
