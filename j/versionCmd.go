package j

import (
  log "github.com/sirupsen/logrus"
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
    log.Info("journalist ", VERSION)
  },
}
