package j

import (
  log "github.com/sirupsen/logrus"
  "net/url"
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/api"
)

var subscribeCmd = &cobra.Command{
  Use:   "subscribe [feed url]",
  Short: "Subscribe to feed",
  Long: "Subscribe to a new feed",
  Args: cobra.MinimumNArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    user := GetApiKey(flagUser, flagPassword)

    feedUrl, err := url.Parse(args[0])
    if err != nil {
      log.Fatal(err)
    }

    var group db.Group
    var grouperr error
    group, grouperr = database.GetGroupByTitleAndUser(flagGroup, user)
    if grouperr != nil {
      log.Info("Group not found, adding ...")
      grouperr = database.AddGroup(&db.Group{
        Title: flagGroup,
        User: user,
      })

      if grouperr != nil {
        log.Fatal(grouperr)
      }

      group, grouperr = database.GetGroupByTitleAndUser(flagGroup, user)
      if grouperr != nil {
        log.Fatal(grouperr)
      }
    }

    err = api.AddOrUpdateFeed(database, feedUrl.String(), group.ID, user)
    if err != nil {
      log.Fatal(err)
    }
    // feed, items, feederr := rss.LoadFeed(feedUrl.String(), group.ID, user)
    // if feederr != nil {
    //   log.Fatal(feederr)
    // }

    // _, upserterr := database.UpsertFeed(&feed, items)
    // if upserterr != nil {
    //   log.Fatal(upserterr)
    // }

    log.Info("Subscribed to ", feedUrl.String())
    return
  },
}

func init() {
  rootCmd.AddCommand(subscribeCmd)
  subscribeCmd.Flags().StringVarP(&flagGroup, "group", "g", "subscriptions", "Group the feed should be assigned to")
  subscribeCmd.Flags().StringVarP(&flagUser, "user", "u", "nobody", "User the feed should be assigned to")
  subscribeCmd.Flags().StringVarP(&flagPassword, "password", "p", "nobody", "Password of user")

  var err error
  database, err = db.InitDatabase()
  if err != nil {
    log.Fatal(err)
  }
}
