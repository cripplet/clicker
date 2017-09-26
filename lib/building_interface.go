package cookie_clicker

var COOKIES_PER_CLICK_LOOKUP float64 = 1

type BuildingInterface interface {
	GetName() string
	GetCost(target int) float64
	GetCPS() float64
}

type BuildingCostFunction func(target int) float64

type Building struct {
	BuildingInterface
	name         string
	costFunction BuildingCostFunction
	cps          float64
}

func NewBuilding(n string, c BuildingCostFunction, cps float64) *Building {
	b := Building{
		name:         n,
		costFunction: c,
		cps:          cps,
	}
	return &b
}

func (self *Building) GetName() string {
	return self.name
}

func (self *Building) GetCPS() float64 {
	return self.cps
}

func (self *Building) GetCost(target int) float64 {
	return self.costFunction(target)
}
