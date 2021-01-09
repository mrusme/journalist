package main

import (
  "os"
  "strconv"
  log "github.com/sirupsen/logrus"
  "github.com/mrusme/journalist/j"
)

func init() {
  log.SetFormatter(&log.JSONFormatter{})
  log.SetOutput(os.Stdout)
  var level log.Level
  levelStr, ok := os.LookupEnv("JOURNALIST_LOG_LEVEL")
  if ok == false {
    level = log.WarnLevel
  } else {
    lvl, levelerr := strconv.ParseUint(levelStr, 10, 32)
    if levelerr != nil {
      level = log.WarnLevel
    } else {
      level = log.Level(lvl)
    }
  }

  log.SetLevel(level)
}

func main() {
  j.Execute()
}
