package env

import (
  "errors"
  "github.com/caarlos0/env/v10"
  "log"
  "os"
)

// LogError log error from env by caarlos0 package. if the exitCode has value it will be used as argument for os.Exit otherwise it will not exit
func LogError(err error, exitCode ...int) {
  var errs env.AggregateError
  if !errors.As(err, &errs) {
    return
  }

  for _, err := range errs.Errors {
    log.Println(err)
  }

  if len(exitCode) > 0 {
    os.Exit(exitCode[0])
  }
}
