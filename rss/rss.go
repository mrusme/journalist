package rss

import (
  // log "github.com/sirupsen/logrus"
  "context"
  "time"
  "strings"
  "strconv"
  "github.com/mmcdole/gofeed"
  // "github.com/mmcdole/gofeed/rss"
  "github.com/mrusme/journalist/db"
)

func LoadFeed(feedUrl string, groupID int64, user string) (db.Feed, []db.Item, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
  defer cancel()

  fp := gofeed.NewParser()
  gfeed, err := fp.ParseURLWithContext(feedUrl, ctx)
  if err != nil {
    return db.Feed{}, []db.Item{}, err
  }

  feedLink := feedUrl
  if gfeed.FeedLink != "" {
    feedLink = gfeed.FeedLink
  }

  feed := db.Feed{
    Title: gfeed.Title,
    Description: gfeed.Description,
    Link: gfeed.Link,
    FeedLink: feedLink,
    Language: gfeed.Language,
    Copyright: gfeed.Copyright,
    Generator: gfeed.Generator,
    Categories: strings.Join(gfeed.Categories, ","),
    Group: groupID,
    User: user,
  }

  if gfeed.Author != nil {
    feed.Author = (*gfeed.Author).Name
  }

  if gfeed.Image != nil {
    feed.Image = (*gfeed.Image).URL
  }

  if gfeed.PublishedParsed != nil {
    feed.CreatedAt = *gfeed.PublishedParsed
  }

  if gfeed.UpdatedParsed != nil {
    feed.UpdatedAt = *gfeed.UpdatedParsed
  }

  var items []db.Item
  for _, gitem := range gfeed.Items {
    item := db.Item{
      Title: gitem.Title,
      Description: gitem.Description,
      Content: gitem.Content,
      Link: gitem.Link,
      Categories: strings.Join(gitem.Categories, ","),
      IsRead: false,
      IsSaved: false,
      User: user,
    }

    if gitem.Author != nil {
      item.Author = (*gitem.Author).Name
    }

    if gitem.Image != nil {
      item.Image = (*gitem.Image).URL
    }

    if gitem.PublishedParsed != nil {
      item.CreatedAt = *gitem.PublishedParsed
    }

    if gitem.UpdatedParsed != nil {
      item.UpdatedAt = *gitem.UpdatedParsed
    }

    item.GUID = GetGUID(gitem.Link + strconv.FormatInt(item.CreatedAt.Unix(), 10))

    items = append(items, item)
  }

  return feed, items, nil
}
