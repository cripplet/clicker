package cookie_clicker

type BuildingType int
type UpgradeID int

const (
	BUILDING_TYPE_MOUSE BuildingType = iota

	BUILDING_TYPE_ENUM_EOF
)

var BUILDING_TYPE_LOOKUP map[BuildingType]string = map[BuildingType]string{
	BUILDING_TYPE_MOUSE: "mouse",
}

// BUILDING_TYPE_REVERSE_LOOKUP is dynamically generated.
var BUILDING_TYPE_REVERSE_LOOKUP map[string]BuildingType = map[string]BuildingType{}

const (
	UPGRADE_ID_REINFORCED_INDEX_FINGER UpgradeID = iota

	UPGRADE_ID_ENUM_EOF
)

var UPGRADE_ID_LOOKUP map[UpgradeID]string = map[UpgradeID]string{
	UPGRADE_ID_REINFORCED_INDEX_FINGER: "reinforced-index-finger",
}

// UPGRADE_ID_REVERSE_LOOKUP is dynamically generated.
var UPGRADE_ID_REVERSE_LOOKUP map[string]UpgradeID = map[string]UpgradeID{}

func init() {
	var b BuildingType
	var u UpgradeID

	for b = BuildingType(0); b < BUILDING_TYPE_ENUM_EOF; b++ {
		if _, present := BUILDINGS_LOOKUP[b]; !present {
			panic("Misconfigured building data")
		}
	}

	for u = UpgradeID(0); u < UPGRADE_ID_ENUM_EOF; u++ {
		if _, present := UPGRADES_LOOKUP[u]; !present {
			panic("Misconfigured upgrade data")
		}
	}

	for buildingType, buildingTypeString := range BUILDING_TYPE_LOOKUP {
		BUILDING_TYPE_REVERSE_LOOKUP[buildingTypeString] = buildingType
	}
	for upgradeID, upgradeIDString := range UPGRADE_ID_LOOKUP {
		UPGRADE_ID_REVERSE_LOOKUP[upgradeIDString] = upgradeID
	}
}
