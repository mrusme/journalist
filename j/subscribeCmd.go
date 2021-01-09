package j

import (
  "os"
  "fmt"
  log "github.com/sirupsen/logrus"
  "net/url"
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/rss"
)

var flagGroup string
var flagUser string
var flagPassword string

var subscribeCmd = &cobra.Command{
  Use:   "subscribe [feed url]",
  Short: "Subscribe to feed",
  Long: "Subscribe to a new feed",
  Args: cobra.MinimumNArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    user := GetApiKey(flagUser, flagPassword)

    feedUrl, err := url.Parse(args[0])
    fmt.Printf("%s\n", feedUrl)
    if err != nil {
      log.Fatal(err)
    }

    var group db.Group
    var grouperr error
    group, grouperr = database.GetGroupByTitleAndUser(flagGroup, user)
    if grouperr != nil {
      log.Println("no group found, adding new one ...")
      grouperr = database.AddGroup(db.Group{
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

    feed, items, feederr := rss.LoadFeed(feedUrl.String(), user)
    if feederr != nil {
      log.Fatal(feederr)
    }

    var feedId int64
    existingFeed, feederr := database.GetFeedByFeedLinkAndUser(feed.FeedLink, user)
    if feederr != nil || existingFeed.ID <= 0 {
      log.Println("subscribing to feed ...")
      feedId, feederr = database.AddFeed(feed, group.ID)

      if feederr != nil {
        log.Fatal(feederr)
      }
    } else {
      feedId = existingFeed.ID
    }

    for _, item := range items {
      _, itemerr := database.AddItem(item, feedId)

      if itemerr != nil {
        log.Debug(itemerr)
      }
    }

    fmt.Printf("%v, %v\n", group.ID, feedId)
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
    fmt.Printf("%+v\n", err)
    os.Exit(1)
  }
}
