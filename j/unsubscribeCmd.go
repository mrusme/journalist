package j

import (
  log "github.com/sirupsen/logrus"
  "net/url"
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
)

var unsubscribeCmd = &cobra.Command{
  Use:   "unsubscribe [feed url]",
  Short: "Unsubscribe from feed",
  Long: "Unsubscribe from a feed",
  Args: cobra.MinimumNArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    user := GetApiKey(flagUser, flagPassword)

    feedUrl, err := url.Parse(args[0])
    if err != nil {
      log.Fatal(err)
    }

    feed, err := database.GetFeedByFeedLinkAndUser(feedUrl.String(), user)
    if err != nil {
      log.Fatal(err)
    }

    err = database.EraseItemsByFeedAndUser(feed.ID, user)
    if err != nil {
      log.Fatal(err)
    }

    err = database.EraseFeedByIDAndUser(feed.ID, user)
    if err != nil {
      log.Fatal(err)
    }

    feedsLeftInGroup, err := database.ListFeedsByGroupAndUser(feed.Group, user)
    if err == nil && len(feedsLeftInGroup) == 0 {
      log.Info("Group empty, removing ...")
      err = database.EraseGroupByID(feed.Group, user)
      log.Debug(err)
    }

    log.Info("Unsubscribed from ", feedUrl.String())
    return
  },
}

func init() {
  rootCmd.AddCommand(unsubscribeCmd)
  unsubscribeCmd.Flags().StringVarP(&flagUser, "user", "u", "nobody", "User the feed should be assigned to")
  unsubscribeCmd.Flags().StringVarP(&flagPassword, "password", "p", "nobody", "Password of user")

  var err error
  database, err = db.InitDatabase()
  if err != nil {
    log.Fatal(err)
  }
}
