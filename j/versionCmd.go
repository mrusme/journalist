package j

import (
  "fmt"

  "github.com/spf13/cobra"
)

var VERSION string

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Display version of journalist",
  Long:  `The version of journalist.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("journalist", VERSION)
  },
}
