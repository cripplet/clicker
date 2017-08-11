package cookie_clicker

func CalculateBuildingUpgradeLoop() {
        var building_type BuildingType
        for _, building_type = range []BuildingType{MOUSE} {
                var upgrade_keys []UpgradeID = BUILDING_UPGRADE_TYPE_REVERSE_LOOKUP[building_type]
                var upgrades []*BuildingUpgrade = []*BuildingUpgrade{}
                for _, upgrade_id := range upgrade_keys {
                        upgrades = append(upgrades, BUILDING_UPGRADE_LIST[upgrade_id])
                }
                var aggregate_upgrade_ratio float64 = GetAggregateUpgradeRatio(upgrades)
                select {
                case BUILDING_UPGRADE_CHANNEL_LOOKUP[building_type] <- aggregate_upgrade_ratio:
                default:
                }
        }
}
