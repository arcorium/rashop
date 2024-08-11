package logger

import (
  sharedConf "github.com/arcorium/rashop/shared/config"
  "github.com/arcorium/rashop/shared/types"
)

var _default = types.Must(NewZapLogger(sharedConf.IsDebug()))

func SetGlobal(logger ILogger) {
  _default = logger
}

func GetGlobal() ILogger {
  return _default
}

type ILogger interface {
  Debugf(format string, args ...any)
  Infof(format string, args ...any)
  Warnf(format string, args ...any)
  Fatalf(format string, args ...any)
  Debug(string)
  Info(string)
  Warn(string)
  Fatal(string)
}

func Debugf(format string, args ...any) {
  _default.Debugf(format, args...)
}

func Infof(format string, args ...any) {
  _default.Infof(format, args...)
}

func Warnf(format string, args ...any) {
  _default.Warnf(format, args...)
}

func Fatalf(format string, args ...any) {
  _default.Fatalf(format, args...)
}

func Debug(s string) {
  _default.Debug(s)
}

func Info(s string) {
  _default.Infof(s)
}

func Warn(s string) {
  _default.Warn(s)
}

func Fatal(s string) {
  _default.Fatalf(s)
}
