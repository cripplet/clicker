package cookie_clicker

import (
	"time"
)

type GameState struct {
	done_channel chan bool
}

// TODO(cripplet): Move global upgrade tables, etc. into GameState.
func MakeGameState() GameState {
	return GameState{
		done_channel: make(chan bool),
	}
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
