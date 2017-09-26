package cookie_clicker

type BuildingType int
type UpgradeID int

const (
	BUILDING_TYPE_MOUSE BuildingType = iota
	BUILDING_TYPE_GRANDMA

	BUILDING_TYPE_ENUM_EOF
)

var BUILDING_TYPE_LOOKUP map[BuildingType]string = map[BuildingType]string{
	BUILDING_TYPE_MOUSE:   "mouse",
	BUILDING_TYPE_GRANDMA: "grandma",
}
var BUILDING_TYPE_REVERSE_LOOKUP map[string]BuildingType = map[string]BuildingType{}

const (
	UPGRADE_ID_REINFORCED_INDEX_FINGER UpgradeID = iota
	UPGRADE_ID_THOUSAND_FINGERS
	UPGRADE_ID_FORWARDS_FROM_GRANDMA

	UPGRADE_ID_ENUM_EOF
)

var UPGRADE_ID_LOOKUP map[UpgradeID]string = map[UpgradeID]string{
	UPGRADE_ID_REINFORCED_INDEX_FINGER: "reinforced-index-finger",
	UPGRADE_ID_THOUSAND_FINGERS:        "thousand-fingers",
	UPGRADE_ID_FORWARDS_FROM_GRANDMA:   "forwards-from-grandma",
}
var UPGRADE_ID_REVERSE_LOOKUP map[string]UpgradeID = map[string]UpgradeID{}

func init() {
	for buildingType, buildingTypeString := range BUILDING_TYPE_LOOKUP {
		BUILDING_TYPE_REVERSE_LOOKUP[buildingTypeString] = buildingType
	}
	for upgradeID, upgradeIDString := range UPGRADE_ID_LOOKUP {
		UPGRADE_ID_REVERSE_LOOKUP[upgradeIDString] = upgradeID
	}
}
