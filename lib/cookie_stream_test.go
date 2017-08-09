package cookie_stream


import (
    "testing"
)


func TestMakeMouseStream(t *testing.T) {
    var m MouseStream = MakeMouseStream()
    if m.name != "Mouse" {
      t.Error("Expected BuildingStream name Mouse, got ", m.name)
    }
    if m.base_cps != 1 {
      t.Error("Expected BuildingStream base_cpc 1, got ", m.base_cps)
    }
}


func TestMouseStreamMine(t *testing.T) {
    var m MouseStream = MakeMouseStream()
    go m.Mine()
    cookies := <- m.cookie_channel
    if cookies <= 0 {
      t.Error("Expected some cookies, got ", cookies)
    }
}


func TestClickStreamMine(t *testing.T) {
    var c ClickStream = MakeClickStream()
    go c.Mine()
    c.Click()
    c.Click()
    cookies := <- c.cookie_channel
    if cookies != 2 {
      t.Error("Expected 2 cookies, got ", cookies)
    }
}


func TestClickStreamMineRace(t *testing.T) {
    var c ClickStream = MakeClickStream()
    go c.Mine()
    var cookies float64
    for i := 0; i < 20; i++ {
      c.Click()
      cookies += <- c.cookie_channel
      c.Click()
    }
    cookies += <- c.cookie_channel
    if cookies != 40 {
      t.Error("Expected 40 cookies, got ", cookies)
    }
}
