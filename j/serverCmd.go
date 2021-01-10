package j

import (
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/api"
  log "github.com/sirupsen/logrus"
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
      log.Fatal(err)
    }

    api.Server(database)
    database.DB.Close()
  },
}
