package config

import (
  "os"
  "strconv"
  "sync"
  "sync/atomic"
)

var once sync.Once
var debug atomic.Bool

func isDebug(env string) bool {
  once.Do(func() {
    s, ok := os.LookupEnv(env)
    if !ok {
      debug.Store(true)
    }
    parseBool, _ := strconv.ParseBool(s)
    debug.Store(!parseBool)
  })
  return debug.Load()
}

func IsDebug() bool {
  return isDebug("ARSHOP_RELEASE")
}
