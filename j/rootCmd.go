package j

import (
  "fmt"
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
  "os"
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
    fmt.Printf("%+v\n", err)
    os.Exit(-1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
}

func initConfig() {
}
