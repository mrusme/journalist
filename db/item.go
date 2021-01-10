package db

import (
  "time"
)

type Item struct {
  ID                int64           `db:"id",json:"id,omitempty"`
  GUID              string          `db:"guid",json:"guid,omitempty"`
  Title             string          `db:"title",json:"title,omitempty"`
  Description       string          `db:"description",json:"description,omitempty"`
  Content           string          `db:"content",json:"content,omitempty"`
  Link              string          `db:"link",json:"link,omitempty"`
  Author            string          `db:"author",json:"author,omitempty"`
  Image             string          `db:"image",json:"image,omitempty"`
  Categories        string          `db:"categories",json:"categories,omitempty"`
  IsRead            bool            `db:"is_read",json:"is_read,omitempty"`
  IsSaved           bool            `db:"is_saved",json:"is_saved,omitempty"`
  Feed              int64           `db:"feed",json:"feed,omitempty"`
  User              string          `db:"user",json:"user,omitempty"`
  CreatedAt         time.Time       `db:"created_at",json:"created_at,omitempty"`
  UpdatedAt         time.Time       `db:"updated_at",json:"updated_at,omitempty"`
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
