package cookie_clicker

import (
	"math"
)

type BuildingCostFunction func(current int) float64

func MouseCost(current int) float64 {
	return math.Pow(2, float64(current))
}

var BUILDING_COST_LOOKUP map[BuildingType]BuildingCostFunction = map[BuildingType]BuildingCostFunction{
	BUILDING_TYPE_MOUSE: MouseCost,
}
