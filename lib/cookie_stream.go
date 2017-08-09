package cookie_stream


import (
    "sync"
    "time"
)


const EPOCH_SLEEP_TIME time.Duration = time.Millisecond * 100


type CookieStreamStruct struct {
  cookie_channel chan float64
  name string
}


type CookieStreamInterface interface {
  Mine()
}


type BuildingStream struct {
  CookieStreamStruct
  CookieStreamInterface
  n_instances int
  base_cps float64
}


type ClickStream struct {
  CookieStreamStruct
  CookieStreamInterface
  base_cpc float64
  clicks int
  clicks_lock sync.Mutex  // TOOD(cripplet): Replace with channel communication instead.
}


type MouseStream struct {
  BuildingStream
}


func MakeCookieStreamStruct(name string) CookieStreamStruct {
  return CookieStreamStruct{
      name: name,
      cookie_channel: make(chan float64),
  }
}


func MakeBuildingStream(name string, base_cps float64) BuildingStream {
  return BuildingStream{
      CookieStreamStruct: MakeCookieStreamStruct(name),
      base_cps: base_cps,
  }
}


func MakeClickStream() ClickStream {
  return ClickStream{
      CookieStreamStruct: MakeCookieStreamStruct("Click"),
      base_cpc: 1,
  }
}


func (self *ClickStream) Click() {
  self.clicks_lock.Lock()
  defer self.clicks_lock.Unlock()
  self.clicks += 1
}


func (self *ClickStream) GetAndResetClicks() int {
  self.clicks_lock.Lock()
  defer self.clicks_lock.Unlock()
  var clicks int
  clicks, self.clicks = self.clicks, 0
  return clicks
}


func (self *ClickStream) Mine() {
  for {
    time.Sleep(EPOCH_SLEEP_TIME)
    self.cookie_channel <- float64(self.GetAndResetClicks()) * self.base_cpc
  }
}


func MakeMouseStream() MouseStream {
  return MouseStream{
      BuildingStream: MakeBuildingStream("Mouse", 1),
  }
}


func (self *BuildingStream) Mine() {
  var last_sent time.Time = time.Now()
  for {
    time.Sleep(EPOCH_SLEEP_TIME)
    var now time.Time = time.Now()
    var cookies float64 = self.base_cps * now.Sub(last_sent).Seconds()
    last_sent = now
    self.cookie_channel <- cookies
  }
}
