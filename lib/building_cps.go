package cookie_clicker

var COOKIES_PER_CLICK_LOOKUP float64 = 1

var BUILDING_CPS_LOOKUP map[BuildingType]float64 = map[BuildingType]float64{
	BUILDING_TYPE_MOUSE:   0.2,
	BUILDING_TYPE_GRANDMA: 1,
}
