package j

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/api"
)

var subscriptionsCmd = &cobra.Command{
  Use:   "subscriptions",
  Short: "List subscriptions",
  Long: "List all subscriptions",
  Run: func(cmd *cobra.Command, args []string) {
    user := api.GetApiKey(flagUser, flagPassword)

    feeds, err := database.ListFeedsByUser(user)
    if err != nil {
      log.Fatal(err)
    }

    for _, feed := range feeds {
      log.Info(feed.FeedLink)
    }

    return
  },
}

func init() {
  rootCmd.AddCommand(subscriptionsCmd)
  subscriptionsCmd.Flags().StringVarP(&flagUser, "user", "u", "nobody", "User to show the subscriptions for")
  subscriptionsCmd.Flags().StringVarP(&flagPassword, "password", "p", "nobody", "Password of user")

  var err error
  database, err = db.InitDatabase()
  if err != nil {
    log.Fatal(err)
  }
}
