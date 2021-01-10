package db

import (
  "time"
  log "github.com/sirupsen/logrus"
)

type Feed struct {
  ID                int64           `db:"id",json:"id,omitempty"`
  Title             string          `db:"title",json:"title,omitempty"`
  Description       string          `db:"description",json:"description,omitempty"`
  Link              string          `db:"link",json:"link,omitempty"`
  FeedLink          string          `db:"feed_link",json:"feed_link,omitempty"`
  Author            string          `db:"author",json:"author,omitempty"`
  Language          string          `db:"language",json:"language,omitempty"`
  Image             string          `db:"image",json:"image,omitempty"`
  Copyright         string          `db:"copyright",json:"copyright,omitempty"`
  Generator         string          `db:"generator",json:"generator,omitempty"`
  Categories        string          `db:"categories",json:"categories,omitempty"`
  Group             int64           `db:"group",json:"group,omitempty"`
  User              string          `db:"user",json:"user,omitempty"`
  CreatedAt         time.Time       `db:"created_at",json:"created_at,omitempty"`
  UpdatedAt         time.Time       `db:"updated_at",json:"update_at,omitempty"`
}

func (database *Database) AddFeed(feed Feed, groupID int64) (int64, error) {
  var id int64
  err := database.DB.QueryRow(`
    INSERT INTO feeds (
      "title",
      "description",
      "link",
      "feed_link",
      "author",
      "language",
      "image",
      "copyright",
      "generator",
      "categories",
      "group",
      "user",
      "created_at",
      "updated_at"
    )
    VALUES (
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7,
      $8,
      $9,
      $10,
      $11,
      $12,
      $13,
      $14
    ) RETURNING "id"
  `,
    feed.Title,
    feed.Description,
    feed.Link,
    feed.FeedLink,
    feed.Author,
    feed.Language,
    feed.Image,
    feed.Copyright,
    feed.Generator,
    feed.Categories,
    groupID,
    feed.User,
    feed.CreatedAt,
    feed.UpdatedAt,
  ).Scan(&id)
  if err != nil {
    return -1, err
  }

  return id, err
}

func (database *Database) GetFeedByFeedLinkAndUser(feedLink string, user string) (Feed, error) {
  var ret Feed

  err := database.DB.Get(&ret, `
    SELECT * FROM feeds WHERE "feed_link" = $1 AND "user" = $2
  `, feedLink, user)

  return ret, err
}

func (database *Database) UpdateFeed(feed Feed) (error) {
  _, err := database.DB.Exec(`
    UPDATE feeds SET
      "title" = $1,
      "description" = $2,
      "link" = $3,
      "feed_link" = $4,
      "author" = $5,
      "language" = $6,
      "image" = $7,
      "copyright" = $8,
      "generator" = $9,
      "categories" = $10,
      "group" = $11,
      "user" = $12,
      "updated_at" = $13
    WHERE "id" = $14
  `,
    feed.Title,
    feed.Description,
    feed.Link,
    feed.FeedLink,
    feed.Author,
    feed.Language,
    feed.Image,
    feed.Copyright,
    feed.Generator,
    feed.Categories,
    feed.Group,
    feed.User,
    feed.UpdatedAt,
    feed.ID,
  )
  return err
}

func (database *Database) UpsertFeed(feed Feed, items []Item) ([]int64, error) {
  var feedID int64
  log.Debug("Checking if feed was already subscribed to ...")
  existingFeed, feederr := database.GetFeedByFeedLinkAndUser(feed.FeedLink, feed.User)
  if feederr != nil || existingFeed.ID <= 0 {
    log.Debug(feederr)
    log.Debug("Subscribing to feed ...")
    log.Debug(feed)
    feedID, feederr = database.AddFeed(feed, feed.Group)

    if feederr != nil {
      return []int64{}, feederr
    }
  } else {
    feedID = existingFeed.ID
    feed.ID = existingFeed.ID
    feed.Group = existingFeed.Group

    log.Debug("Already subscribed to feed, updating ...")
    updateerr := database.UpdateFeed(feed)
    if updateerr != nil {
      log.Debug(updateerr)
    }
  }
  log.Debug("Feed ID: ", feedID)
  log.Debug("Refreshing items ...")

  var itemIDs []int64
  for _, item := range items {
    itemID, itemerr := database.AddItem(item, feedID)

    if itemerr != nil {
      existingItem, geterr := database.GetItemByGUIDAndUser(item.GUID, item.User)
      if geterr != nil {
        log.Debug(geterr)
      } else {
        item.ID = existingItem.ID
        item.Feed = existingItem.Feed
        if item.UpdatedAt.After(existingItem.UpdatedAt) {
          item.IsRead = false
        } else {
          item.IsRead = existingItem.IsRead
        }
        item.IsSaved = existingItem.IsSaved
        updateerr := database.UpdateItem(item)
        if updateerr != nil {
          log.Debug(updateerr)
        }
      }
      log.Debug(itemerr)
    } else {
      log.Debug("Added new item with ID:", itemID)
      itemIDs = append(itemIDs, itemID)
    }
  }

  return itemIDs, nil
}

// func (database *Database) EraseFeed(feed Feed) (error) {
// }

func (database *Database) ListFeeds() ([]Feed, error) {
  var ret []Feed

  res, err := database.DB.Queryx(`
    SELECT * FROM feeds
  `)

  if err != nil {
    return ret, err
  }

  for res.Next() {
    var feed Feed

    err := res.StructScan(&feed)
    if err != nil {
      return ret, err
    }

    ret = append(ret, feed)
  }

  return ret, err
}

func (database *Database) ListFeedsByUser(user string) ([]Feed, error) {
  var ret []Feed

  res, err := database.DB.Queryx(`
    SELECT * FROM feeds WHERE "user" = $1
  `, user)

  if err != nil {
    return ret, err
  }

  for res.Next() {
    var feed Feed

    err := res.StructScan(&feed)
    if err != nil {
      return ret, err
    }

    ret = append(ret, feed)
  }

  return ret, err
}

func (database *Database) ListFeedsByGroupAndUser(groupID int64, user string) ([]Feed, error) {
  var ret []Feed

  res, err := database.DB.Queryx(`
    SELECT * FROM feeds
    WHERE
      "group" = $1
    AND
      "user" = $2
  `,
    groupID,
    user,
  )

  if err != nil {
    return ret, err
  }

  for res.Next() {
    var feed Feed

    err := res.StructScan(&feed)
    if err != nil {
      return ret, err
    }

    ret = append(ret, feed)
  }

  return ret, err
}
