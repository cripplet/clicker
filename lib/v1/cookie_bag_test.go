package cookie_clicker


import (
    "sync"
    "testing"
)


func TestMakeCookieBag(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  if c.cookies != 0 {
    t.Error("Expected empty cookie bag, got ", c.cookies)
  }
}


func TestCookieBagPeekBag(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  c.cookies = 100
  var actual float64 = c.PeekBag()
  if actual != 100 {
    t.Error("Expected 100 cookies in bag, got ", actual)
  }
}


func TestCookieBagPutBag(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  if !c.PutBag(100) {
    t.Error("Could not put 100 cookies in bag.")
  }

  var actual float64 = c.PeekBag()
  if actual != 100 {
    t.Error("Expected 100 cookies in bag, got ", actual)
  }
}


func TestCookieBagPutBagNegative(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  if c.PutBag(-1) {
    t.Error("Put negative cookies in bag.")
  }

  var actual float64 = c.PeekBag()
  if actual != 0 {
    t.Error("Expected empty cookie bag, got ", actual)
  }
}


func TestCookieBagGrabBag(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  c.PutBag(100)

  if !c.GrabBag(100) {
    t.Error("Cound not grab 100 cookies from bag.")
  }

  var actual float64 = c.PeekBag()
  if actual != 0 {
    t.Error("Expected empty cookie bag, got ", actual)
  }
}


func TestCookieBagGrabBagNegative(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  c.PutBag(100)

  if c.GrabBag(-100) {
    t.Error("Grabbed a negative amount of cookies from the bag.")
  }

  var actual float64 = c.PeekBag()
  if actual != 100 {
    t.Error("Expected 100 cookies in bag, got ", actual)
  }
}


func TestCookieBagGrabBagTooExpensive(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  c.PutBag(100)

  if c.GrabBag(101) {
    t.Error("Grabbed too many cookies from the bag.")
  }

  var actual float64 = c.PeekBag()
  if actual != 100 {
    t.Error("Expected 100 cookies in bag, got ", actual)
  }
}


func PutBagWorker(wg *sync.WaitGroup, b *CookieBag, n float64) {
  b.PutBag(n)
  wg.Done()
}


func TestCookieBagPutBagRace(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  var wg sync.WaitGroup = sync.WaitGroup{}

  for i := 0; i < 10; i++ {
    wg.Add(1)
    go PutBagWorker(&wg, &c, 100)
  }
  wg.Wait()

  var actual float64 = c.PeekBag()
  if actual != 1000 {
    t.Error("Expected 1000 cookies in bag, got ", actual)
  }
}


func GrabBagWorker(wg *sync.WaitGroup, b *CookieBag, n float64) {
  b.GrabBag(n)
  wg.Done()
}


func TestCookieBagGrabBagRace(t *testing.T) {
  var c CookieBag = MakeCookieBag()
  var wg sync.WaitGroup = sync.WaitGroup{}
  c.PutBag(1000)

  for i := 0; i < 10; i++ {
    wg.Add(1)
    go GrabBagWorker(&wg, &c, 100)
  }
  wg.Wait()

  var actual float64 = c.PeekBag()
  if actual != 0 {
    t.Error("Expected empty cookie bag, got ", actual)
  }
}
