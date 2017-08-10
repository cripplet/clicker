package cookie_clicker


import (
//    "reflect"
)


type StreamType int


const (
    CLICK StreamType = iota
    MOUSE
)


type GameState struct {
  stream_list map[StreamType]*CookieStreamInterface // BaseStreamStruct
  click_stream *ClickStream
}


// TODO(cripplet): Add configuration file as input.
func MakeGameState() GameState {
  var click_stream ClickStream = MakeClickStream()
  var mouse_stream MouseStream = MakeMouseStream()

  var stream_list map[StreamType]*CookieStreamInterface = map[StreamType]*CookieStreamInterface{
      CLICK: &(click_stream.CookieStreamInterface),
      MOUSE: &(mouse_stream.CookieStreamInterface),
  }

  return GameState{
      stream_list: stream_list,
      click_stream: &click_stream,
  }
}

/*
func (self *GameState) CookieBagDaemon() {
  var cookie_stream_switch []reflect.SelectCase = make([]reflect.SelectCase, len(self.stream_list))
  for k, stream := range self.stream_list {
    cookie_stream_switch[k] = reflect.SelectCase{
        Dir: reflect.SelectRecv,
        Chan: reflect.ValueOf((*stream).GetCookieChannel()),
    }
  }

  for {
    chosen, value, ok := reflect.Select(cookie_stream_switch)
    println(chosen, value.Float(), ok)
  }
}
 */

func (self *GameState) Start() {
  for _, stream := range self.stream_list {
    go stream.Mine()
  }
  // go self.CookieBagDaemon()
}
