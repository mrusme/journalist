package j

import (
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/api"
  "fmt"
  "os"
)

func init() {
  rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
  Use:   "server",
  Short: "Run journalist server",
  Long:  `Run the journalist server.`,
  Run: func(cmd *cobra.Command, args []string) {
    var err error
    database, err = db.InitDatabase()
    if err != nil {
      fmt.Printf("%+v\n", err)
      os.Exit(1)
    }

    api.Server()
  },
}
