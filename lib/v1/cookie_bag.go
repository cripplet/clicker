package cookie_clicker


import (
    "sync"
)


type CookieBag struct {
  cookies float64
  cookies_lock sync.Mutex
}


func MakeCookieBag() CookieBag {
  return CookieBag{}
}


func (self *CookieBag) PeekBag() float64 {
  self.cookies_lock.Lock()
  defer self.cookies_lock.Unlock()
  return self.cookies
}


func (self *CookieBag) PutBag(n float64) bool {
  if n <= 0 {
    return false
  }

  self.cookies_lock.Lock()
  defer self.cookies_lock.Unlock()

  self.cookies += n
  return true
}


func (self *CookieBag) GrabBag(n float64) bool {
  if n <= 0 {
    return false
  }

  self.cookies_lock.Lock()
  defer self.cookies_lock.Unlock()

  if self.cookies < n {
    return false
  }

  self.cookies -= n
  return true
}
