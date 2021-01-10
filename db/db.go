package db

import (
  "os"
  "strings"
  "regexp"
  "time"
  "errors"
  log "github.com/sirupsen/logrus"

  _ "database/sql"
  "github.com/jmoiron/sqlx"
  _ "github.com/jackc/pgx/v4/stdlib"
)

var schema = `
CREATE TABLE IF NOT EXISTS groups (
    "id" SERIAL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "title_unix" TEXT NOT NULL,
    "user" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "groups_title_unix" ON "groups"("title_unix","user");

CREATE TABLE IF NOT EXISTS feeds (
    "id" SERIAL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "link" TEXT NOT NULL,
    "feed_link" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "language" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "copyright" TEXT NOT NULL,
    "generator" TEXT NOT NULL,
    "categories" TEXT NOT NULL,
    "group" INT NOT NULL,
    "user" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT fk_group FOREIGN KEY("group") REFERENCES groups("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "feeds_feed_link" ON "feeds"("feed_link","user");

CREATE TABLE IF NOT EXISTS items (
    "id" SERIAL PRIMARY KEY,
    "guid" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "link" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "categories" TEXT NOT NULL,
    "is_read" BOOL NOT NULL,
    "is_saved" BOOL NOT NULL,
    "feed" INT NOT NULL,
    "user" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT fk_feed FOREIGN KEY("feed") REFERENCES feeds("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "items_guid" ON "items"("guid","user");
`

type Database struct {
  DB *sqlx.DB
}

func InitDatabase() (*Database, error) {
  dbconnection, ok := os.LookupEnv("JOURNALIST_DB")
  if ok == false || dbconnection == "" {
    return nil, errors.New("please `export JOURNALIST_DB` with the database connection string, e.g. 'postgres://user:secret@localhost:5432/journalist?sslmode=disable'")
  }

  db, err := sqlx.Open("pgx", dbconnection)
  if err != nil {
    return nil, err
  }

  err = db.Ping()
  if err != nil {
    db.Close()
    return nil, err
  }

  db.MustExec(schema)

  database := Database{db}
  return &database, nil
}

func (database *Database) AddGroup(group Group) (error) {
  _, err := database.DB.Exec(`
    INSERT INTO groups ("title", "title_unix", "user", "created_at", "updated_at")
    VALUES ($1, $2, $3, $4, $5)
  `, group.Title, GetUnixName(group.Title), group.User, time.Now(), time.Now())
  return err
}

func (database *Database) GetGroupByID(groupID uint) (Group, error) {
  var ret Group

  err := database.DB.Get(&ret, `
    SELECT * FROM groups WHERE "id" = $1
  `, groupID)

  if err != nil {
    return ret, err
  }

  return ret, err
}

func (database *Database) GetGroupByTitleAndUser(title string, user string) (Group, error) {
  var ret Group

  err := database.DB.Get(&ret, `
    SELECT * FROM groups WHERE "title_unix" = $1 AND "user" = $2
  `, GetUnixName(title), user)

  return ret, err
}

func (database *Database) UpdateGroup(group Group) (error) {
  _, err := database.DB.Exec(`
    UPDATE groups SET ? WHERE "id" = ?
  `, &group, group.ID)
  return err
}

func (database *Database) EraseGroup(group Group) (error) {
  _, err := database.DB.Exec(`
    DELETE FROM groups WHERE "id" = ?
  `, group.ID)
  return err
}

func (database *Database) ListGroupsByUser(user string) ([]Group, error) {
  var ret []Group

  res, err := database.DB.Queryx(`
    SELECT * FROM groups WHERE "user" = $1
  `, user)

  if err != nil {
    return ret, err
  }

  for res.Next() {
    var group Group

    err := res.StructScan(&group)
    if err != nil {
      return ret, err
    }

    ret = append(ret, group)
  }

  return ret, err
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
  existingFeed, feederr := database.GetFeedByFeedLinkAndUser(feed.FeedLink, feed.User)
  if feederr != nil || existingFeed.ID <= 0 {
    log.Debug(feederr)
    log.Debug("Subscribing to feed ...")
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

func (database *Database) AddItem(item Item, feedId int64) (int64, error) {
  var id int64
  err := database.DB.QueryRow(`
    INSERT INTO items (
      "guid",
      "title",
      "description",
      "content",
      "link",
      "author",
      "image",
      "categories",
      "is_read",
      "is_saved",
      "feed",
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
    item.GUID,
    item.Title,
    item.Description,
    item.Content,
    item.Link,
    item.Author,
    item.Image,
    item.Categories,
    item.IsRead,
    item.IsSaved,
    feedId,
    item.User,
    item.CreatedAt,
    item.UpdatedAt,
  ).Scan(&id)
  if err != nil {
    return -1, err
  }

  return id, err
}

func (database *Database) GetItemByGUIDAndUser(itemGUID string, user string) (Item, error) {
  var ret Item

  err := database.DB.Get(&ret, `
    SELECT * FROM items WHERE
      "guid" = $1
    AND
      "user" = $2
  `,
    itemGUID,
    user,
  )

  if err != nil {
    return ret, err
  }

  return ret, err
}

func (database *Database) UpdateItem(item Item) (error) {
  _, err := database.DB.Exec(`
    UPDATE items SET
      "guid" = $1,
      "title" = $2,
      "description" = $3,
      "content" = $4,
      "link" = $5,
      "author" = $6,
      "image" = $7,
      "categories" = $8,
      "is_read" = $9,
      "is_saved" = $10,
      "feed" = $11,
      "user" = $12,
      "updated_at" = $13
    WHERE "id" = $14
  `,
    item.GUID,
    item.Title,
    item.Description,
    item.Content,
    item.Link,
    item.Author,
    item.Image,
    item.Categories,
    item.IsRead,
    item.IsSaved,
    item.Feed,
    item.User,
    item.UpdatedAt,
    item.ID,
  )
  return err
}

func (database *Database) UpdateItemByIDAsRead(itemID int64, read bool, user string) (error) {
  _, err := database.DB.Exec(`
    UPDATE items SET
      "is_read" = $2
    WHERE
      "id" = $1
    AND
      "user" = $3
  `,
    itemID,
    read,
    user,
  )
  return err
}

func (database *Database) UpdateItemByIDAsSaved(itemID int64, saved bool, user string) (error) {
  _, err := database.DB.Exec(`
    UPDATE items SET
      "is_saved" = $2
    WHERE
      "id" = $1
    AND
      "user" = $3
  `,
    itemID,
    saved,
    user,
  )
  return err
}

func (database *Database) UpdateItemsByBeforeAsRead(before time.Time, read bool, user string) (error) {
  _, err := database.DB.Exec(`
    UPDATE items SET
      "is_read" = $2
    WHERE
      "created_at" < $1
    AND
      "user" = $3
  `,
    before,
    read,
    user,
  )
  return err
}

func (database *Database) UpdateItemsByGroupAsRead(groupID int64, read bool, user string) (error) {
  feeds, err := database.ListFeedsByGroupAndUser(groupID, user)
  if err != nil {
    return err
  }

  var feedIDs []int64
  for _, feed := range feeds {
    feedIDs = append(feedIDs, feed.ID)
  }

  query, args, err := sqlx.In(`
    UPDATE items SET
      "is_read" = ?
    WHERE
      "group" IN (?)
    AND
      "user" = ?
  `,
    read,
    feedIDs,
    user,
  )
  if err != nil {
    return err
  }

  _, err = database.DB.Exec(
    database.DB.Rebind(query), args...
  )

  return err
}


// func (database *Database) EraseItem(item Item) (error) {
// }

func (database *Database) ListItemsByUser(user string, sinceID int64) ([]Item, error) {
  var ret []Item

  res, err := database.DB.Queryx(`
    SELECT * FROM items
    WHERE
      "user" = $1
    AND
      "id" >= $2
  `,
    user,
    sinceID,
  )

  if err != nil {
    return ret, err
  }

  for res.Next() {
    var item Item

    err := res.StructScan(&item)
    if err != nil {
      return ret, err
    }

    ret = append(ret, item)
  }

  return ret, err
}

func (database *Database) ListItemsByIDsAndUser(ids []int64, user string) ([]Item, error) {
  var ret []Item

  query, args, err := sqlx.In(`
    SELECT * FROM items
    WHERE
      "id" IN (?)
    AND
      "user" = ?
  `,
    ids,
    user,
  )
  if err != nil {
    return []Item{}, err
  }

  res, err := database.DB.Queryx(
    database.DB.Rebind(query), args...
  )

  if err != nil {
    return ret, err
  }

  for res.Next() {
    var item Item

    err := res.StructScan(&item)
    if err != nil {
      return ret, err
    }

    ret = append(ret, item)
  }

  return ret, err
}

func (database *Database) ListUnreadItemsByUser(user string, sinceID int64) ([]Item, error) {
  var ret []Item

  res, err := database.DB.Queryx(`
    SELECT * FROM items
    WHERE
      "user" = $1
    AND
      "id" >= $2
    AND
      "is_read" = FALSE
  `,
    user,
    sinceID,
  )

  if err != nil {
    return ret, err
  }

  for res.Next() {
    var item Item

    err := res.StructScan(&item)
    if err != nil {
      return ret, err
    }

    ret = append(ret, item)
  }

  return ret, err
}

func GetUnixName(name string) string {
  reg, regerr := regexp.Compile("[^a-zA-Z0-9]+")
  if regerr != nil {
      return ""
  }

  id := strings.ToLower(reg.ReplaceAllString(name, ""))

  return id
}
