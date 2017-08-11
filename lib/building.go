package cookie_clicker

import (
	"time"
)

const EPOCH_TIME time.Duration = time.Millisecond * 500
const CHANNEL_TIMEOUT time.Duration = time.Nanosecond

var BUILDING_CPS_LOOKUP map[BuildingType]float64 = map[BuildingType]float64{
	MOUSE:   0.1,
	GRANDMA: 1,
	FARM:    5,
}

var BUILDING_UPGRADE_CHANNEL_LOOKUP map[BuildingType]chan float64 = map[BuildingType]chan float64{
	MOUSE:   make(chan float64, 1),
	GRANDMA: make(chan float64, 1),
	FARM:    make(chan float64, 1),
}

type CookieStream struct {
	building_type BuildingType

	cookie_channel      chan float64
	cookie_done_channel chan bool

	upgrade_ratio float64
}

func MakeCookieStream(t BuildingType) CookieStream {
	return CookieStream{
		building_type:       t,
		cookie_channel:      make(chan float64, 1),
		cookie_done_channel: make(chan bool),
		upgrade_ratio:       1,
	}
}

func generateCookieLoop(c *CookieStream, l time.Time, n time.Time) bool {
	select {
	case c.upgrade_ratio = <-BUILDING_UPGRADE_CHANNEL_LOOKUP[c.building_type]:
	case <-time.After(CHANNEL_TIMEOUT):
	}
	select {
	case c.cookie_channel <- c.upgrade_ratio * BUILDING_CPS_LOOKUP[c.building_type] * n.Sub(l).Seconds():
		return true
	default:
		return false
	}
}

func generateCookie(c *CookieStream) {
	var l, n time.Time
	l = time.Now()
	for {
		n = time.Now()
		if generateCookieLoop(c, l, n) {
			l = n
		}
		select {
		case <-c.cookie_done_channel:
			return
		case <-time.After(CHANNEL_TIMEOUT):
		}
	}
}

func StartCookieStream(c *CookieStream) {
	go generateCookie(c)
}

func StopCookieStream(c *CookieStream) {
	c.cookie_done_channel <- true
}

func GetCookie(c *CookieStream) float64 {
	return <-c.cookie_channel
}
