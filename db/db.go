package db

import (
  "os"
  "sort"
  "strings"
  "strconv"
  "regexp"
  "log"
  // "fmt"
  "errors"
  "encoding/json"
  "github.com/tidwall/buntdb"
  "github.com/google/uuid"
)

type Database struct {
  DB *buntdb.DB
}

func InitDatabase() (*Database, error) {
  dbfile, ok := os.LookupEnv("JOURNALIST_DB")
  if ok == false || dbfile == "" {
    return nil, errors.New("please `export JOURNALIST_DB` to the location the geld database should be stored at")
  }

  db, err := buntdb.Open(dbfile)
  if err != nil {
    return nil, err
  }

  db.CreateIndex("group", "*", buntdb.IndexJSON("group"))

  database := Database{db}
  return &database, nil
}

func (database *Database) NewID() (string) {
  id, err := uuid.NewRandom()
  if err != nil {
    log.Fatalln("could not generate UUID: %+v", err)
  }
  return id.String()
}

func (database *Database) GetIncrementID(user string, entity string) (int, error) {
  var idInt int

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    id, dberr := tx.Get(user + ":id:" + entity, false)
    if dberr != nil {
      return nil
    }

    if idInt, dberr = strconv.Atoi(id); dberr != nil {
      idInt = 0
    }

    idInt++

    _, _, seterr := tx.Set(user + ":id:" + entity, strconv.Itoa(idInt), nil)
    if seterr != nil {
      return seterr
    }

    return nil
  })

  return idInt, dberr
}

func (database *Database) AddItem(user string, item Item) (string, error) {
  id := database.NewID()

  incID, incerr := database.GetIncrementID(user, "item")
  if incerr != nil {
    return "", incerr
  }
  item.IncID = incID

  itemJson, jsonerr := json.Marshal(item)
  if jsonerr != nil {
    return id, jsonerr
  }

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, _, seterr := tx.Set(user + ":item:" + id, string(itemJson), nil)
    if seterr != nil {
      return seterr
    }

    return nil
  })

  return id, dberr
}

func (database *Database) GetItem(user string, itemId string) (Item, error) {
  var item Item

  dberr := database.DB.View(func(tx *buntdb.Tx) error {
    value, err := tx.Get(user + ":item:" + itemId, false)
    if err != nil {
      return nil
    }

    json.Unmarshal([]byte(value), &item)

    return nil
  })

  return item, dberr
}

func (database *Database) UpdateItem(user string, item Item) (string, error) {
  itemJson, jsonerr := json.Marshal(item)
  if jsonerr != nil {
    return item.ID, jsonerr
  }

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, _, seerr := tx.Set(user + ":item:" + item.ID, string(itemJson), nil)
    if seerr != nil {
      return seerr
    }

    return nil
  })

  return item.ID, dberr
}

func (database *Database) EraseItem(user string, id string) (error) {
  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, delerr := tx.Delete(user + ":item:" + id)
    if delerr != nil {
      return delerr
    }

    return nil
  })

  return dberr
}

func (database *Database) ListItems(user string) ([]Item, error) {
  var items []Item

  dberr := database.DB.View(func(tx *buntdb.Tx) error {
    tx.AscendKeys(user + ":item:*", func(key, value string) bool {
      var item Item
      json.Unmarshal([]byte(value), &item)

      item.SetIDFromDatabaseKey(key)

      items = append(items, item)
      return true
    })

    return nil
  })

  sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
  return items, dberr
}

func (database *Database) AddGroup(user string, group Group) (string, error) {
  id := database.NewID()

  incID, incerr := database.GetIncrementID(user, "group")
  if incerr != nil {
    return "", incerr
  }
  group.IncID = incID

  groupJson, jsonerr := json.Marshal(group)
  if jsonerr != nil {
    return id, jsonerr
  }

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, _, seterr := tx.Set(user + ":group:" + id, string(groupJson), nil)
    if seterr != nil {
      return seterr
    }

    return nil
  })

  return id, dberr
}

func (database *Database) ListGroups(user string) ([]Group, error) {
  var groups []Group

  dberr := database.DB.View(func(tx *buntdb.Tx) error {
    tx.AscendKeys(user + ":group:*", func(key, value string) bool {
      var group Group
      json.Unmarshal([]byte(value), &group)

      group.SetIDFromDatabaseKey(key)

      groups = append(groups, group)
      return true
    })

    return nil
  })

  sort.Slice(groups, func(i, j int) bool { return groups[i].IncID < groups[j].IncID })
  return groups, dberr
}

func (database *Database) UpdateGroup(user string, groupName string, group Group) (error) {
  groupJson, jsonerr := json.Marshal(group)
  if jsonerr != nil {
    return jsonerr
  }

  foundGroup, founderr := database.GetGroup(user, groupName)
  if founderr != nil {
    return founderr
  }

  groupId := foundGroup.ID

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, _, sperr := tx.Set(user + ":group:" + groupId, string(groupJson), nil)
    if sperr != nil {
      return sperr
    }

    return nil
  })

  return dberr
}

func (database *Database) GetGroup(user string, groupName string) (Group, error) {
  var group Group
  found := false
  groupUnixName := GetUnixName(groupName)

  groups, dberr := database.ListGroups(user)
  if dberr != nil {
    return group, dberr
  }

  for _, g := range groups {
    if GetUnixName(g.Title) == groupUnixName {
      group = g
      found = true
      break
    }
  }

  if found == false {
    return group, errors.New("No group found")
  }

  return group, dberr
}

func GetIDFromDatabaseKey(key string) (string, error) {
  splitKey := strings.Split(key, ":")

  if len(splitKey) < 3 || len(splitKey) > 3 {
    return "", errors.New("not a valid database key")
  }

  return splitKey[2], nil
}

func GetUnixName(name string) string {
  reg, regerr := regexp.Compile("[^a-zA-Z0-9]+")
  if regerr != nil {
      return ""
  }

  id := strings.ToLower(reg.ReplaceAllString(name, ""))

  return id
}
