package cookie_clicker

import (
	"time"
)

const EPOCH_TIME time.Duration = time.Millisecond * 500

var BUILDING_CPS_LOOKUP map[BuildingType]float64 = map[BuildingType]float64{
	MOUSE:   0.1,
	GRANDMA: 1,
	FARM:    5,
}

var BUILDING_UPGRADE_CHANNEL_LOOKUP map[BuildingType]chan float64 = map[BuildingType]chan float64{
	MOUSE:   make(chan float64),
	GRANDMA: make(chan float64),
	FARM:    make(chan float64),
}

type CookieStream struct {
	building_type BuildingType

	cookie_channel chan float64
	done_channel   chan bool

	upgrade_ratio float64
}

func MakeCookieStream(t BuildingType) CookieStream {
	return CookieStream{
		building_type:  t,
		cookie_channel: make(chan float64),
		done_channel:   make(chan bool),
		upgrade_ratio:  1,
	}
}

func generateCookieLoop(c *CookieStream, l time.Time, n time.Time) bool {
	select {
	case c.upgrade_ratio = <-BUILDING_UPGRADE_CHANNEL_LOOKUP[c.building_type]:
	case <-time.After(time.Nanosecond):
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
		case <-c.done_channel:
			return
		case <-time.After(time.Nanosecond):
		}
	}
}
