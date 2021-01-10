package j

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
)

var database *db.Database

var flagGroup string
var flagUser string
var flagPassword string

var rootCmd = &cobra.Command{
  Use:   "journalist",
  Short: "RSS aggregator",
  Long:  `A RSS aggregator.`,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    log.Fatal(err)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
}

func initConfig() {
}
