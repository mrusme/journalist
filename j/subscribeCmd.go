package j

import (
  "os"
  "fmt"
  "log"
  "net/url"
  "github.com/spf13/cobra"
  "github.com/mrusme/journalist/db"
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
    if err != nil {
      log.Fatal(err)
    }

    // _, err = database.AddEntry(user, newEntry, isRunning)
    // if err != nil {
    //   fmt.Printf("%s %+v\n", CharError, err)
    //   os.Exit(1)
    // }

    // fmt.Printf(newEntry.GetOutputForTrack(isRunning, false))
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
