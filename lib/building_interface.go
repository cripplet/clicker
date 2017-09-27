package cookie_clicker

var COOKIES_PER_CLICK_LOOKUP float64 = 1

// BuildingInterface is the public API for all buildings.
type BuildingInterface interface {
	GetName() string
	GetCost(target int) float64
	GetDescription() string
	GetCPS() float64
}

type buildingCostFunction func(target int) float64

type standardBuilding struct {
	BuildingInterface
	name         string
	description  string
	costFunction buildingCostFunction
	cps          float64
}

func newStandardBuilding(n string, d string, c buildingCostFunction, cps float64) *standardBuilding {
	b := standardBuilding{
		name:         n,
		description:  d,
		costFunction: c,
		cps:          cps,
	}
	return &b
}

func (self *standardBuilding) GetName() string {
	return self.name
}

func (self *standardBuilding) GetDescription() string {
	return self.description
}

func (self *standardBuilding) GetCPS() float64 {
	return self.cps
}

func (self *standardBuilding) GetCost(target int) float64 {
	return self.costFunction(target)
}
