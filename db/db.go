package db

import (
  "os"
  "strings"
  "regexp"
  // "fmt"
  "time"
  "errors"
  // "encoding/json"

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

// func (database *Database) AddFeed(feed Feed) (string, error) {
// }

func (database *Database) GetFeedByLinkAndUser(link string, user string) (Feed, error) {
  var ret Feed

  err := database.DB.Get(&ret, `
    SELECT * FROM feeds WHERE "link" = $1 AND "user" = $2
  `, link, user)

  return ret, err
}

// func (database *Database) UpdateFeed(feed Feed) (string, error) {
// }

// func (database *Database) EraseFeed(feed Feed) (error) {
// }

// func (database *Database) ListFeedsByUser(user string) ([]Feed, error) {
// }

// func (database *Database) AddItem(item Item) (string, error) {
// }

// func (database *Database) GetItem(item Item) (Item, error) {
// }

// func (database *Database) UpdateItem(item Item) (string, error) {
// }

// func (database *Database) EraseItem(item Item) (error) {
// }

// func (database *Database) ListItemsByUser(user string) ([]Item, error) {
// }

func GetUnixName(name string) string {
  reg, regerr := regexp.Compile("[^a-zA-Z0-9]+")
  if regerr != nil {
      return ""
  }

  id := strings.ToLower(reg.ReplaceAllString(name, ""))

  return id
}
