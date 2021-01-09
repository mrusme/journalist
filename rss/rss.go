package rss

import (
  // "fmt"
  "context"
  "time"
  "github.com/mmcdole/gofeed"
  // "github.com/mmcdole/gofeed/rss"
)

func LoadFeed(feedUrl string) (*gofeed.Feed) {
  ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
  defer cancel()

  fp := gofeed.NewParser()
  feed, _ := fp.ParseURLWithContext(feedUrl, ctx)

  return feed
}
