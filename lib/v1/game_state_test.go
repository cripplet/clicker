package cookie_clicker


import (
    "testing"
)


func TestRunGameState(t *testing.T) {
  var g GameState = MakeGameState()
  go g.Start()
  for {}
}
