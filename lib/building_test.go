package cookie_clicker

import (
	"fmt"
	"math"
	"testing"
	"time"
)

const ONE_SECOND_DURATION time.Duration = time.Duration(1) * time.Second
const EPSILON_PERCENT float64 = 0.05

func TestMakeCookieStream(t *testing.T) {
	var s CookieStream = MakeCookieStream(MOUSE)
	if s.building_type != MOUSE {
		t.Error(fmt.Sprintf("Expected building type %d, got %d", MOUSE, s.building_type))
	}

	if s.upgrade_ratio != 1 {
		t.Error(fmt.Sprintf("Expected upgrade ratio of %e, got %e", 1, s.upgrade_ratio))
	}
}

func TestGenerateCookieLoop(t *testing.T) {
	var s CookieStream = MakeCookieStream(MOUSE)

	var n time.Time = time.Now()
	go generateCookieLoop(&s, n, n.Add(ONE_SECOND_DURATION))

	var n_cookies float64 = GetCookie(&s)
	if n_cookies != BUILDING_CPS_LOOKUP[MOUSE] {
		t.Error(fmt.Sprintf("Expected %e cookies, got %e", BUILDING_CPS_LOOKUP[MOUSE], n_cookies))
	}
}

func TestGenerateUpgrade(t *testing.T) {
	const upgrade_ratio float64 = 2

	var s CookieStream = MakeCookieStream(MOUSE)
	go func() {
		BUILDING_UPGRADE_CHANNEL_LOOKUP[MOUSE] <- upgrade_ratio
	}()

	time.Sleep(EPOCH_TIME) // Wait until we're "sure" the upgrade channel is blocking.

	var n time.Time = time.Now()
	go generateCookieLoop(&s, n, n.Add(ONE_SECOND_DURATION))

	var n_cookies float64 = GetCookie(&s)
	if n_cookies != upgrade_ratio*BUILDING_CPS_LOOKUP[MOUSE] {
		t.Error(fmt.Sprintf("Expected %e cookies, got %e", upgrade_ratio*BUILDING_CPS_LOOKUP[MOUSE], n_cookies))
	}
}

func TestGenerateCookie(t *testing.T) {
	var s CookieStream = MakeCookieStream(MOUSE)
	go generateCookie(&s)
	time.Sleep(time.Second)

	var n_cookies float64
	var current_buffer_len int = len(s.cookie_channel)
	for i := 0; i < current_buffer_len+1; i++ {
		n_cookies += GetCookie(&s)
	}

	if math.Abs(n_cookies-BUILDING_CPS_LOOKUP[MOUSE]) > math.Max(n_cookies, BUILDING_CPS_LOOKUP[MOUSE])*EPSILON_PERCENT {
		t.Error(fmt.Sprintf("Expected %e cookies, got %e", BUILDING_CPS_LOOKUP[MOUSE], n_cookies))
	}
	s.cookie_done_channel <- true
}
