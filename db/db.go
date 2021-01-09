package db

import (
  "os"
  "strings"
  "regexp"
  "fmt"
  "time"
  "errors"
  // "encoding/json"
  "github.com/genjidb/genji"
  "github.com/genjidb/genji/document"
)

type Database struct {
  DB *genji.DB
}

func InitDatabase() (*Database, error) {
  dbfile, ok := os.LookupEnv("JOURNALIST_DB")
  if ok == false || dbfile == "" {
    return nil, errors.New("please `export JOURNALIST_DB` to the location the geld database should be stored at")
  }

  db, err := genji.Open(dbfile)
  if err != nil {
    return nil, err
  }

  err = db.Exec("CREATE TABLE groups (id INTEGER PRIMARY KEY, title TEXT NOT NULL, user TEXT NOT NULL, createdAt INTEGER, updatedAt INTEGER)")
  err = db.Exec("CREATE TABLE feeds")
  err = db.Exec("CREATE TABLE items")

  database := Database{db}
  return &database, nil
}

func (database *Database) AddGroup(group Group) (error) {
  err := database.DB.Exec(`
    INSERT INTO groups (title, user, createdAt, updatedAt)
    VALUES (?, ?, ?, ?)
  `, group.Title, group.User, time.Now().Unix(), time.Now().Unix())
  return err
}

func (database *Database) GetGroupByID(groupID uint) (Group, error) {
  var ret Group

  d, err := database.DB.QueryDocument(`
    SELECT * FROM groups WHERE id = ?
  `, &groupID)

  if err != nil {
    return ret, err
  }

  err = document.StructScan(d, &ret)

  return ret, err
}

func (database *Database) GetGroupByTitleAndUser(title string, user string) (Group, error) {
  var ret Group
  retFound := false

  groups, err := database.ListGroupsByUser(user)
  if err != nil {
    return ret, err
  }

  unixTitle := GetUnixName(title)

  for _, group := range groups {
    if GetUnixName(group.Title) == unixTitle {
      fmt.Printf("Found group! %v\n", group.ID)
      ret = group
      retFound = true
      break
    }
  }

  if retFound == false {
    return ret, errors.New("Not found")
  }

  return ret, nil
}

func (database *Database) UpdateGroup(group Group) (error) {
  err := database.DB.Exec(`
    UPDATE groups SET ? WHERE id = ?
  `, &group, group.ID)
  return err
}

func (database *Database) EraseGroup(group Group) (error) {
  err := database.DB.Exec(`
    DELETE FROM groups WHERE id = ?
  `, group.ID)
  return err
}

func (database *Database) ListGroupsByUser(user string) ([]Group, error) {
  var ret []Group

  res, err := database.DB.Query(`
    SELECT * FROM groups WHERE user = ?
  `, &user)
  defer res.Close()

  if err != nil {
    return ret, err
  }

  err = res.Iterate(func(d document.Document) (error) {
    var group Group
    err = document.StructScan(d, &group)
    if err != nil {
      return err
    }

    ret = append(ret, group)
    return nil
  })

  return ret, err
}

// func (database *Database) AddFeed(feed Feed) (string, error) {
// }

// func (database *Database) GetFeed(feed Feed) (Feed, error) {
// }

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
