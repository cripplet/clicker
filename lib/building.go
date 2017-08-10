package cookie_clicker


import (
    "time"
)


const EPOCH_TIME time.Duration = time.Millisecond * 500


type CommandType int
type BuildingType int


const (
    SYN CommandType = iota
)


const (
    MOUSE BuildingType = iota
    GRANDMA
    FARM
)


var BUILDING_CPS_LOOKUP map[BuildingType]float64 = map[BuildingType]float64{
    MOUSE: 0.1,
    GRANDMA: 1,
    FARM: 5,
}


var BUILDING_UPGRADE_CHANNEL_LOOKUP map[BuildingType]chan float64 = map[BuildingType]chan float64{
    MOUSE: make(chan float64),
    GRANDMA: make(chan float64),
    FARM: make(chan float64),
}


type CookieStream struct {
  building_type BuildingType

  command_channel chan CommandType
  cookie_channel chan float64

  upgrade_ratio float64
}


func MakeCookieStream(t BuildingType) CookieStream {
  return CookieStream{
      building_type: t,
      command_channel: make(chan CommandType),
      cookie_channel: make(chan float64),
      upgrade_ratio: 1,
  }
}


func dispatchCookieStreamCommandChannel(c *CookieStream, command CommandType) {
}


func generateCookieLoop(c *CookieStream, l time.Time, n time.Time) {
  select {
    case c.upgrade_ratio = <- BUILDING_UPGRADE_CHANNEL_LOOKUP[c.building_type]:
    default:
  }
  c.cookie_channel <- c.upgrade_ratio * BUILDING_CPS_LOOKUP[c.building_type] * n.Sub(l).Seconds()
  select {
    case command := <- c.command_channel:
      dispatchCookieStreamCommandChannel(c, command)
    default:
  }
}


func generateCookie(c *CookieStream) {
  var l, n time.Time
  for {
    n = time.Now()
    generateCookieLoop(c, l, n)
    l = n
  }
}
