package cookie_stream


import(
    "sync"
)


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
  clicks_lock sync.Mutex
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


func (self *ClickStream) Mine() {
  for {
    var clicks int
    {
      self.clicks_lock.Lock()
      defer self.clicks_lock.Unlock()
      clicks = self.clicks
      self.clicks = 0
    }
    self.cookie_channel <- float64(clicks) * self.base_cpc
  }
}

func MakeMouseStream() MouseStream {
  return MouseStream{
      BuildingStream: MakeBuildingStream("Mouse", 1),
  }
}


func (self *BuildingStream) Mine() {
  for {
    self.cookie_channel <- self.base_cps
  }
}
