package db

import (
  "time"

  _ "database/sql"
  "github.com/jmoiron/sqlx"
  _ "github.com/jackc/pgx/v4/stdlib"
  readability "github.com/go-shiori/go-readability"
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
  ReadableTitle     string          `db:"readable_title",json:"readable_title,omitempty"`
  ReadableAuthor    string          `db:"readable_author",json:"readable_author,omitempty"`
  ReadableExcerpt   string          `db:"readable_excerpt",json:"readable_excerpt,omitempty"`
  ReadableSiteName  string          `db:"readable_site_name",json:"readable_site_name,omitempty"`
  ReadableImage     string          `db:"readable_image",json:"readable_image,omitempty"`
  ReadableContent   string          `db:"readable_content",json:"readable_content,omitempty"`
  ReadableText      string          `db:"readable_text",json:"readable_text,omitempty"`
  IsRead            bool            `db:"is_read",json:"is_read,omitempty"`
  IsSaved           bool            `db:"is_saved",json:"is_saved,omitempty"`
  Feed              int64           `db:"feed",json:"feed,omitempty"`
  User              string          `db:"user",json:"user,omitempty"`
  CreatedAt         time.Time       `db:"created_at",json:"created_at,omitempty"`
  UpdatedAt         time.Time       `db:"updated_at",json:"updated_at,omitempty"`
}

func (item *Item) AssignReadableFromArticle(article *readability.Article) {
  item.ReadableTitle = article.Title
  item.ReadableAuthor = article.Byline
  item.ReadableExcerpt = article.Excerpt
  item.ReadableSiteName = article.SiteName
  item.ReadableImage = article.Image
  item.ReadableContent = article.Content
  item.ReadableText = article.TextContent
}

func (database *Database) AddItem(item *Item, feedId int64) (int64, error) {
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
      "readable_title",
      "readable_author",
      "readable_excerpt",
      "readable_site_name",
      "readable_image",
      "readable_content",
      "readable_text",
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
      $14,
      $15,
      $16,
      $17,
      $18,
      $19,
      $20,
      $21
    )
    ON CONFLICT ("guid", "user") DO NOTHING
    RETURNING "id"
  `,
    item.GUID,
    item.Title,
    item.Description,
    item.Content,
    item.Link,
    item.Author,
    item.Image,
    item.Categories,
    item.ReadableTitle,
    item.ReadableAuthor,
    item.ReadableExcerpt,
    item.ReadableSiteName,
    item.ReadableImage,
    item.ReadableContent,
    item.ReadableText,
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

func (database *Database) UpdateItem(item *Item) (error) {
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
      "readable_title" = $9,
      "readable_author" = $10,
      "readable_excerpt" = $11,
      "readable_site_name" = $12,
      "readable_image" = $13,
      "readable_content" = $14,
      "readable_text" = $15,
      "is_read" = $16,
      "is_saved" = $17,
      "feed" = $18,
      "user" = $19,
      "updated_at" = $20
    WHERE "id" = $21
  `,
    item.GUID,
    item.Title,
    item.Description,
    item.Content,
    item.Link,
    item.Author,
    item.Image,
    item.Categories,
    item.ReadableTitle,
    item.ReadableAuthor,
    item.ReadableExcerpt,
    item.ReadableSiteName,
    item.ReadableImage,
    item.ReadableContent,
    item.ReadableText,
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

func (database *Database) UpdateItemsByFeedAndBeforeAsRead(feedID int64, before time.Time, read bool, user string) (error) {
  _, err := database.DB.Exec(`
    UPDATE items SET
      "is_read" = $3
    WHERE
      "feed" = $1
    AND
      "created_at" < $2
    AND
      "user" = $4
  `,
    feedID,
    before,
    read,
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


func (database *Database) EraseItemsByFeedAndUser(feedID int64, user string) (error) {
  _, err := database.DB.Exec(`
    DELETE FROM items
    WHERE
      "feed" = $1
    AND
      "user" = $2
  `,
    feedID,
    user,
  )
  return err
}

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
