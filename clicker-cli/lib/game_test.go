package cc_cli_lib

import (
	"github.com/cripplet/clicker/lib"
	"testing"
	"time"
)

var start time.Time = time.Now()
var end time.Time = start.Add(time.Second)

func TestGameStateStartStop(t *testing.T) {
	s := cookie_clicker.NewGameState()
	g := NewGame(s)
	g.game.Load(cookie_clicker.GameStateData{})
	if g.game.GetCPS(start, end) != 0 {
		t.Errorf("Unexpected CPS: %e", g.game.GetCPS(start, end))
	}

	var i float64
	for i = 0; i < g.game.GetBuildingCost()[cookie_clicker.BUILDING_TYPE_MOUSE](0); i += g.game.GetCookiesPerClick() {
		g.game.Click()
	}

	if !g.game.BuyBuilding(cookie_clicker.BUILDING_TYPE_MOUSE) {
		t.Errorf("Could not buy building")
	}

	current_cookies := g.game.GetCookies()

	if g.game.GetCPS(start, end) == 0 {
		t.Errorf("Game CPS not altered")
	}

	g.Start()
	time.Sleep(2 * time.Second)
	g.Stop()

	if (g.game.GetCookies() - current_cookies) < s.GetCPS(start, end) {
		t.Errorf("Expected more than %e cookies, got %e", g.game.GetCPS(start, end), g.game.GetCookies()-current_cookies)
	}
}
