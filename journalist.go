package main

import (
  "os"
  log "github.com/sirupsen/logrus"
  "github.com/mrusme/journalist/j"
  "github.com/mrusme/journalist/common"
)

func init() {
  log.SetFormatter(&log.JSONFormatter{})
  log.SetOutput(os.Stdout)

  level := common.LookupIntEnv("JOURNALIST_LOG_LEVEL", int64(log.WarnLevel))
  log.SetLevel(log.Level(level))
}

func main() {
  j.Execute()
}
