package cookie_clicker

import (
	"time"
)

type UpgradeID int

const (
	REINFORCED_INDEX_FINGER UpgradeID = iota
)

var BUILDING_UPGRADE_TYPE_LOOKUP map[UpgradeID]BuildingType = map[UpgradeID]BuildingType{
	REINFORCED_INDEX_FINGER: MOUSE,
}

var BUILDING_UPGRADE_TYPE_REVERSE_LOOKUP map[BuildingType][]UpgradeID = map[BuildingType][]UpgradeID{
	MOUSE: []UpgradeID{
		REINFORCED_INDEX_FINGER,
	},
}

type BuildingUpgrade struct {
	name          string
	upgrade_ratio float64
	activated_channel chan bool
	activated_done_channel chan bool
	activated_send_channel chan bool
	activated_send_done_channel chan bool
	cost float64
}

func activateBuildingUpgradeLoop(b *BuildingUpgrade) {
	select {
	case b.activated_send_channel <- true:
	default:
	}
}

func activateBuildingUpgrade(b *BuildingUpgrade) {
	for {
		activateBuildingUpgradeLoop(b)
		select {
		case <-b.activated_send_done_channel:
			return
		case <-time.After(CHANNEL_TIMEOUT):
		}
	}
}

func listenBuildingUpgradeLoop(b *BuildingUpgrade) {
	select {
	case <-b.activated_channel:
		activateBuildingUpgrade(b)
	// TODO(cripplet): <-activated_cancel_channel
	case <-time.After(CHANNEL_TIMEOUT):
	}
}

func listenBuildingUpgrade(b *BuildingUpgrade) {
	for {
		listenBuildingUpgradeLoop(b)
		select {
		case <-b.activated_done_channel:
			return
		case <-time.After(CHANNEL_TIMEOUT):
		}
	}
}

func StartBuildingUpgrade(b *BuildingUpgrade) {
        go listenBuildingUpgrade(b)
}

func StopBuildingUpgrade(b *BuildingUpgrade) {
	b.activated_send_done_channel <- true
	b.activated_done_channel <- true
}

func ActivateBuildingUpgrade(b *BuildingUpgrade) {
	b.activated_channel <- true
}

func MakeBuildingUpgrade(name string, upgrade_ratio float64, cost float64) BuildingUpgrade {
	return BuildingUpgrade{
		name: name,
		upgrade_ratio: upgrade_ratio,
		activated_channel: make(chan bool, 1),
		activated_done_channel: make(chan bool),
		activated_send_channel: make(chan bool, 1),
		activated_send_done_channel: make(chan bool),
		cost: cost,
	}
}

var REINFORCED_INDEX_FINGER_UPGRADE BuildingUpgrade = MakeBuildingUpgrade("Reinforced Index Finger", 2, 100)

var BUILDING_UPGRADE_LIST map[UpgradeID]*BuildingUpgrade = map[UpgradeID]*BuildingUpgrade{
	REINFORCED_INDEX_FINGER: &REINFORCED_INDEX_FINGER_UPGRADE,
}

func ActivateUpgrade(u UpgradeID) {
	go ActivateBuildingUpgrade(BUILDING_UPGRADE_LIST[u])
}

func GetAggregateUpgradeRatio(upgrades []*BuildingUpgrade) float64 {
	var aggregate_upgrade_ratio float64 = 1
	for _, upgrade := range upgrades {
		var is_activated bool = false
		select {
		case is_activated = <-upgrade.activated_send_channel:
		case<-time.After(CHANNEL_TIMEOUT):
		}
		if is_activated {
			aggregate_upgrade_ratio *= upgrade.upgrade_ratio
		}
	}
	return aggregate_upgrade_ratio
}
