package j

import (
  "fmt"
  "github.com/spf13/cobra"
  "os"
)

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
