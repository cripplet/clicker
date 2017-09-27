package cookie_clicker

var COOKIES_PER_CLICK_LOOKUP float64 = 1

type BuildingInterface interface {
	GetName() string
	GetCost(target int) float64
	GetCPS() float64
}

type buildingCostFunction func(target int) float64

type standardBuilding struct {
	BuildingInterface
	name         string
	costFunction buildingCostFunction
	cps          float64
}

func newStandardBuilding(n string, c buildingCostFunction, cps float64) *standardBuilding {
	b := standardBuilding{
		name:         n,
		costFunction: c,
		cps:          cps,
	}
	return &b
}

func (self *standardBuilding) GetName() string {
	return self.name
}

func (self *standardBuilding) GetCPS() float64 {
	return self.cps
}

func (self *standardBuilding) GetCost(target int) float64 {
	return self.costFunction(target)
}
