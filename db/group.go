package db

import (
  "time"
)

type Group struct {
  ID                int64           `db:"id",json:"id,omitempty"`
  Title             string          `db:"title",json:"title,omitempty"`
  TitleUnix         string          `db:"title_unix",json:"title_unix,omitempty"`
  User              string          `db:"user",json:"user,omitempty"`
  CreatedAt         time.Time       `db:"created_at",json:"created_at,omitempty"`
  UpdatedAt         time.Time       `db:"updated_at",json:"updated_at,omitempty"`
}

func (database *Database) AddGroup(group *Group) (error) {
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

func (database *Database) UpdateGroup(group *Group) (error) {
  _, err := database.DB.Exec(`
    UPDATE groups SET ? WHERE "id" = ?
  `, &group, group.ID)
  return err
}

func (database *Database) EraseGroupByID(groupID int64, user string) (error) {
  _, err := database.DB.Exec(`
    DELETE FROM groups
    WHERE
      "id" = $1
    AND
      "user" = $2
  `,
    groupID,
    user,
  )
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
