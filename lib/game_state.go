package cookie_clicker

import (
	"time"
)

var BUILDING_TYPES []BuildingType = []BuildingType{
	MOUSE,
}

type GameState struct {
	n_buildings_channels map[BuildingType]chan int
	n_buildings_done_channels map[BuildingType]chan bool
	n_buildings_broadcast_channels map[BuildingType] chan int
	n_buildings_broadcast_done_channels map[BuildingType]chan bool
	done_channel chan bool
	n_buildings int
}

// TODO(cripplet): Move global upgrade tables, etc. into GameState.
func MakeGameState() GameState {
	var g GameState = GameState{
		done_channel: make(chan bool),
	}
	for _, building_type := range BUILDING_TYPES {
		g.n_buildings_channels[building_type] = make(chan int, 1)
		g.n_buildings_done_channels[building_type] = make(chan bool)
		g.n_buildings_broadcast_channels[building_type] = make(chan int, 1)
		g.n_buildings_broadcast_done_channels[building_type] = make(chan bool)
	}
	return g
}

func listenBuildingChannelsLoop(g* GameState, b BuildingType, n int) int {
	var n_buildings int = n
	select {
	case n_buildings = <-g.n_buildings_channels[b]:
	case <-time.After(EPOCH_TIME):
	}
	return n_buildings
}

func listenBuildingChannels(g* GameState, b BuildingType) {
	var n_buildings int
	for {
		n_buildings = startBuildingChannelsLoop(g, b, n_buildings)
		select {
		case <-g.n_buildings_done_channels[b]:
			return
		case <-time.After(CHANNEL_TIMEOUT):
		}
	}
}

func startBuildingChannels(g *GameState) {
	for _, building_type :=  range BUILDING_TYPES {
		go listenBuildingChannels(b, building_type)
	}
}

func stopBuildingChannels(g* GameState, b BuildingType) {
	g.n_buildings_done_channels[b] <- true
}

func GenerateGameLoop(g *GameState) {
	var b BuildingType
	for _, b = range []BuildingType{
		MOUSE,
	} {
		CalculateBuildingUpgrade(b)
	}
}

func GenerateGame(g *GameState) {
	for {
		GenerateGameLoop(g)
		select {
		case <-g.done_channel:
			return
		case <-time.After(EPOCH_TIME):
		}
	}
}

func StartGame(g *GameState) {
	go GenerateGame(g)
}

func StopGame(g *GameState) {
	g.done_channel <- true
}

func CalculateBuildingUpgrade(b BuildingType) {
	var upgrade_keys []UpgradeID = BUILDING_UPGRADE_TYPE_REVERSE_LOOKUP[b]
	var upgrades []*BuildingUpgrade = []*BuildingUpgrade{}
	for _, upgrade_id := range upgrade_keys {
		upgrades = append(upgrades, BUILDING_UPGRADE_LIST[upgrade_id])
	}
	var aggregate_upgrade_ratio float64 = GetAggregateUpgradeRatio(upgrades)
	select {
	case BUILDING_UPGRADE_CHANNEL_LOOKUP[b] <- aggregate_upgrade_ratio:
	default:
	}
}
