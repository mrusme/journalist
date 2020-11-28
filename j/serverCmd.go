package j

import (
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/api"
)

func init() {
  rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
  Use:   "server",
  Short: "Run journalist server",
  Long:  `Run the journalist server.`,
  Run: func(cmd *cobra.Command, args []string) {
    api.Server()
  },
}
