package cc_cli_lib

import (
	"github.com/cripplet/clicker/lib"
	"time"
)

var EPOCH_MILLISECONDS time.Duration = time.Duration(time.Millisecond * 100)

type Game struct {
	game               *cookie_clicker.GameStateStruct
	mainLoopDoneSignal chan bool
}

func NewGame(g *cookie_clicker.GameStateStruct) *Game {
	ng := Game{
		game: g,
	}
	return &ng
}

func (self *Game) Start() {
	self.mainLoopDoneSignal = make(chan bool)
	go self.loopDaemon()
}

func (self *Game) Stop() {
	self.mainLoopDoneSignal <- true
}

func (self *Game) loopDaemon() {
	t := time.NewTicker(EPOCH_MILLISECONDS)
	last_updated := time.Now()
	for {
		select {
		case <-t.C:
			current_time := time.Now()
			self.game.MineCookies(last_updated, current_time)
			last_updated = current_time
			continue
		case <-self.mainLoopDoneSignal:
			t.Stop()
			return
		}
	}
}
