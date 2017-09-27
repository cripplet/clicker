package cookie_clicker

import (
	"math"
)

var BUILDINGS_LOOKUP map[BuildingType]BuildingInterface = map[BuildingType]BuildingInterface{
	BUILDING_TYPE_MOUSE: newStandardBuilding(
		"Mouse",
		"Autoclicks once every 10 seconds",
		func(target int) float64 { return math.Pow(2, float64(target)) + 15. },
		0.2,
	),
}
