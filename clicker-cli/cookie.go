package main

import (
    "fmt"
    "time"
    "github.com/cripplet/clicker/somedep"
)

type Building interface {
  GetBase() float64
  GetName() string
  GetCost() int
}

type BuildingStruct struct {
  base float64
  name string
  cost int
  n int
}

func (b BuildingStruct) GetBase() float64 {
  return b.base
}

func (b BuildingStruct) GetName() string {
  return b.name
}

func (b BuildingStruct) GetCost() int {
  return b.cost
}

type Mouse struct {
  BuildingStruct
}

func ConstructMouse() Mouse {
  return Mouse{
      BuildingStruct: BuildingStruct{
          base: 1,
          name: "mouse",
          cost: 10,
      },
  }
}

type Modifier interface {
  GetName() string
  GetModifier() float64
}

type CookiePerSecond struct {
  buildings []*Building
  modifiers []*Modifier
  last_update_time time.Time
  cookie_stream chan float64
}

type CookieBag struct {
  current float64
}

func ConstructCookieBag() CookieBag {
  return CookieBag{}
}

func (c *CookieBag) GetCurrent() int {
  return int(c.current)
}

func (c *CookieBag) Buy(b Building, cps *CookiePerSecond) bool {
  if c.GrabCookies(float64(b.GetCost())) {
    cps.AddBuilding(&b)
    return true
  }
  return false
}

func (c *CookieBag) GrabCookies(n float64) bool {
  if c.current >= n {
    c.current -= n
    return true
  }
  return false
}

func (c *CookieBag) PutCookies(n float64) {
  c.current += n
}

func ConstructCookiePerSecond() CookiePerSecond {
  return CookiePerSecond{
      last_update_time: time.Now(),
      cookie_stream: make(chan float64),
  }
}

func (c *CookiePerSecond) Click() {
  c.cookie_stream <- 1
}

func (c *CookiePerSecond) AddBuilding(b *Building) {
  c.buildings = append(c.buildings, b)
}

func (c *CookiePerSecond) PeekCookiePerSecond() float64 {
  var base float64 = 0
  for _, b := range c.buildings {
    base += float64((*b).GetBase())
  }

  var aggregate_modifier float64 = 1
  for _, m := range c.modifiers {
    aggregate_modifier *= (*m).GetModifier()
  }

  return base * aggregate_modifier
}

func (c *CookiePerSecond) GetNextCookies() float64 {
  var base float64 = c.PeekCookiePerSecond()

  var n time.Time = time.Now()
  d := n.Sub(c.last_update_time)
  c.last_update_time = n

  return base * d.Seconds()
}

func (c *CookiePerSecond) CookieMiner() {
  for i :=  0; ; i++ {
    c.cookie_stream <- c.GetNextCookies()
    time.Sleep(time.Millisecond * 500)
  }
}

func (cb *CookieBag) CookieReceiver(ch chan float64) {
  for i := 0; ; i++ {
    cb.PutCookies(<- ch)
  }
}

func Play(cb *CookieBag, c *CookiePerSecond) {
  fmt.Println("CPS", "Bank")
  for i := 0; ; i++ {
    fmt.Println(c.PeekCookiePerSecond(), cb.GetCurrent())
    c.Click()
    cb.Buy(ConstructMouse(), c)
    time.Sleep(time.Millisecond * 500)
  }
}

func main() {
  var c CookiePerSecond = ConstructCookiePerSecond()
  var cb CookieBag = ConstructCookieBag()

  go c.CookieMiner()
  go cb.CookieReceiver(c.cookie_stream)

  somedep.Foo()
  Play(&cb, &c)
}
